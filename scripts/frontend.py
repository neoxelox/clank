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

    prettier_files = '"**/*.{ts,svelte,html,scss}"' if file == "." else file
    eslint_files = file
    stylelint_files = '"**/*.{svelte,html,scss}"' if file == "." else file

    with context.cd("frontend"):
        context.run(f"{Tools.Yarn} run types")
        context.run(f"{Tools.Yarn} run svelte-check --tsconfig ./tsconfig.json")
        context.run(f"{Tools.Yarn} run prettier --ignore-path .gitignore --check {prettier_files}")
        context.run(
            f'{Tools.Yarn} run eslint --ext ".ts,.svelte,.html" --max-warnings 0 --ignore-path .gitignore {eslint_files}'
        )
        context.run(f"{Tools.Yarn} run stylelint --ignore-path .gitignore {stylelint_files}")


@task(
    help={
        "file": "[<FILE_PATH>]. If empty, it will format all files.",
    },
)
def format(context, file="."):
    """Run formatter."""

    prettier_files = '"**/*.{ts,svelte,html,scss}"' if file == "." else file
    eslint_files = file
    stylelint_files = '"**/*.{svelte,html,scss}"' if file == "." else file

    with context.cd("frontend"):
        context.run(f"{Tools.Yarn} run types")
        context.run(f"{Tools.Yarn} run prettier --ignore-path .gitignore --write {prettier_files}")
        context.run(
            f'{Tools.Yarn} run eslint --ext ".ts,.svelte,.html" --max-warnings 0 --ignore-path .gitignore --fix {eslint_files}'
        )
        context.run(f"{Tools.Yarn} run stylelint --ignore-path .gitignore --fix {stylelint_files}")


@task(
    help={
        "test": "[<FILE_PATH>]::[<TEST_NAME>]. If empty, it will run all tests.",
    },
)
def test(context, test="."):
    """Run tests."""

    with context.cd("frontend"):
        context.run(f"{Tools.Yarn} run types")
        context.run("echo skipping... no test runner!")


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

    with context.cd("frontend"):
        context.run(
            f"{Tools.Docker} build "
            f"--secret id=env,src=../infra/{target}/.env "
            f"--file ../infra/{target}/frontend.Dockerfile "
            # f"-t frontend-{target}:{context.commit()} "
            f"-t frontend-{target}:latest "
            f"."
        )


@task(
    help={
        "target": "Target environment. If empty, current.",
    }
)
def deploy(context, target=str(Envs.Current)):
    """Deploy to Cloudflare Pages."""

    if Envs.Current != Envs.Ci:
        context.fail(f"deploy command only available in {Envs.Ci} environment!")

    ensure_envfile(context, f"infra/{target}/.env")

    context.run(
        f"{Tools.Docker} run --rm "
        f"--env-file infra/{target}/.env "
        f"frontend-{target}:latest wrangler pages deploy",
        pty=False,
    )
