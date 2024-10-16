import os
import re

from superinvoke import console, rich, task

from .common import ensure_envfile
from .envs import Envs
from .tools import Tools


@task(
    help={
        "file": "[<FILE_PATH>]. If empty, it will lint all files.",
    },
)
def lint(context, file="./..."):
    """Run linter."""

    with context.cd("backend"):
        context.run(f"{Tools.Go} vet {file}")
        context.run(f"{Tools.GolangCILint} run {file} -c .golangci.yaml")
        context.run(f"{Tools.Squawk} -c .squawk.toml migrations/*.sql")


@task(
    help={
        "file": "[<FILE_PATH>]. If empty, it will format all files.",
    },
)
def format(context, file="./..."):
    """Run formatter."""

    with context.cd("backend"):
        context.run(f"{Tools.Go} fmt {file}")
        context.run(f"{Tools.GolangCILint} run {file} -c .golangci.yaml --fix")


@task(
    help={
        "test": "[<PACKAGE_PATH>]::[<TEST_NAME>]. If empty, it will run all tests.",
        "verbose": "Show stdout of tests.",
        "show": "Show coverprofile page.",
    },
)
def test(context, test="", verbose=False, show=False):
    """Run tests."""

    test_arg = "./..."
    if test:
        test = test.split("::")
        if len(test) == 1 and test[0]:
            test_arg = f"{test[0]}/..."
        if len(test) == 2 and test[1]:
            test_arg += f" -run {test[1]}"

    verbose_arg = ""
    if verbose:
        verbose_arg = "-v"

    parallel_arg = ""
    if os.cpu_count():
        parallel_arg = f"--parallel={os.cpu_count()}"

    coverprofile_arg = ""
    if show:
        coverprofile_arg = "-coverprofile=coverage.out"

    with context.cd("backend"):
        result = context.run(
            f"{Tools.GoTestSum} --format=testname --no-color=False -- {verbose_arg} {parallel_arg} -race -count=1 -cover {coverprofile_arg} {test_arg}",
        )

    if "DONE 0 tests" not in result.stdout:
        packages = 0
        coverage = 0.0

        for cover in re.findall(r"[0-9]+\.[0-9]+(?=%)", result.stdout):
            packages += 1
            coverage += float(cover)

        if packages:
            coverage = round(coverage / packages, 1)

        console.print(
            rich.panel.Panel(
                f"Total Coverage ([bold]{packages} pkg[/bold]): [bold green]{coverage}%[/bold green]",
                expand=False,
            )
        )

    if show:
        with context.cd("backend"):
            context.run(f"{Tools.Go} tool cover -html=coverage.out")
            context.remove("coverage.out")


@task(
    help={
        "name": "Migration name.",
    }
)
def migrate(context, name):
    """Create a migration."""

    with context.cd("backend"):
        context.run(f"{Tools.GolangMigrate} create -ext sql -dir migrations -seq -digits 4 {name}")


@task(variadic=True)
def run(context, args):
    """Execute a command."""

    if Envs.Current == Envs.Dev:
        context.run(f"{Tools.Compose} exec -T cli cli {args}", pty=False)
    else:
        context.run(f"{Tools.Compose} run --no-deps --rm cli {args}", pty=False)


@task(
    help={
        "target": "Target environment. If empty, current.",
    }
)
def build(context, target=str(Envs.Current)):
    """Build images."""

    if Envs.Current != Envs.Ci:
        context.fail(f"build command only available in {Envs.Ci} environment!")

    ensure_envfile(context, f"infra/{target}/.env")

    with context.cd("backend"):
        for service in ["api", "worker", "cli"]:
            context.run(
                f"{Tools.Docker} build --file ../infra/{target}/{service}.Dockerfile "
                # f"-t {service}-{target}:{context.commit()} "
                f"-t {service}-{target}:latest "
                f"."
            )


@task(
    help={
        "target": "Target environment. If empty, current.",
    }
)
def deploy(context, target=str(Envs.Current), server=""):
    """Push images to the registry."""

    if Envs.Current != Envs.Ci:
        context.fail(f"deploy command only available in {Envs.Ci} environment!")

    ensure_envfile(context, f"infra/{target}/.env")

    # TODO: Activate this when ready
    # for service in ["api", "worker", "cli"]:
    #     context.run(f"{Tools.Docker} push --all-tags {service}-{target}")

    # TODO: Remove this when ready
    if not server:
        context.fail("currently sshing to a server is needed for deploy!")

    for service in ["api", "worker", "cli"]:
        context.run(f"{Tools.Docker} save -o {service}-{target}-latest.tar {service}-{target}:latest")
        context.run(f"scp {service}-{target}-latest.tar {server}:~/workspace")
        context.remove(f"{service}-{target}-latest.tar")
        context.run(f"ssh -n {server} 'docker load -i ~/workspace/{service}-{target}-latest.tar'")
        context.run(f"ssh -n {server} 'rm -f ~/workspace/{service}-{target}-latest.tar'")

    context.run(f"scp infra/{target}/.env {server}:~/workspace")
    context.run(
        f"ssh -n {server} 'docker compose "
        f"--env-file ~/workspace/.env "
        f"--file ~/workspace/docker-compose.yaml "
        f"up --detach --force-recreate --no-deps --timeout=60 --wait --wait-timeout=60 "
        f"api worker cli'"
    )
    context.run(f"ssh -n {server} 'docker rmi $(docker images --filter dangling=true -q --no-trunc)'")
