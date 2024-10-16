from superinvoke import task

from .common import ensure_envfile
from .envs import Envs
from .tools import Tools


@task(
    help={
        "file": "[<FILE_PATH>]. If empty, it will lint all files.",
    },
)
def lint(context, file="."):
    """Run linter."""

    with context.cd("engine"):
        context.run(f"{Tools.Poetry} run ruff check {file}")
        # TODO: Activate this when ready
        # context.run(f"{Tools.Poetry} run mypy {file}")


@task(
    help={
        "file": "[<FILE_PATH>]. If empty, it will format all files.",
    },
)
def format(context, file="."):
    """Run formatter."""

    with context.cd("engine"):
        context.run(f"{Tools.Poetry} run ruff check --fix-only {file}")
        context.run(f"{Tools.Poetry} run ruff format {file}")


@task(
    help={
        "test": "[<FILE_PATH>]::[<TEST_NAME>]. If empty, it will run all tests.",
    },
)
def test(context, test="."):
    """Run tests."""

    with context.cd("engine"):
        context.run(f"{Tools.Poetry} run pytest -p no:warnings -n auto {test}")


@task(
    help={
        "target": "Target environment. If empty, current.",
    }
)
def build(context, target=str(Envs.Current)):
    """Build image."""

    if Envs.Current != Envs.Ci:
        context.fail(f"build command only available in {Envs.Ci} environment!")

    ensure_envfile(context, f"infra/{target}/.env")

    with context.cd("engine"):
        context.run(
            f"{Tools.Docker} build --file ../infra/{target}/engine.Dockerfile "
            # f"-t engine-{target}:{context.commit()} "
            f"-t engine-{target}:latest "
            f"."
        )


@task(
    help={
        "target": "Target environment. If empty, current.",
    }
)
def deploy(context, target=str(Envs.Current), server=""):
    """Push image to the registry."""

    if Envs.Current != Envs.Ci:
        context.fail(f"deploy command only available in {Envs.Ci} environment!")

    ensure_envfile(context, f"infra/{target}/.env")

    # TODO: Activate this when ready
    # context.run(f"{Tools.Docker} push --all-tags engine-{target}")

    # TODO: Remove this when ready
    if not server:
        context.fail("currently sshing to a server is needed for deploy!")

    context.run(f"{Tools.Docker} save -o engine-{target}-latest.tar engine-{target}:latest")
    context.run(f"scp engine-{target}-latest.tar {server}:~/workspace")
    context.remove(f"engine-{target}-latest.tar")
    context.run(f"ssh -n {server} 'docker load -i ~/workspace/engine-{target}-latest.tar'")
    context.run(f"ssh -n {server} 'rm -f ~/workspace/engine-{target}-latest.tar'")
    context.run(f"scp infra/{target}/.env {server}:~/workspace")
    context.run(
        f"ssh -n {server} 'docker compose "
        f"--env-file ~/workspace/.env "
        f"--file ~/workspace/docker-compose.yaml "
        f"up --detach --force-recreate --no-deps --timeout=60 --wait --wait-timeout=60 "
        f"engine'"
    )
    context.run(f"ssh -n {server} 'docker rmi $(docker images --filter dangling=true -q --no-trunc)'")
