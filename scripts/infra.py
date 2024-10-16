from superinvoke import task

from .tools import Tools


@task(default=True)
def list(context):
    """List containers."""

    context.run(f"{Tools.Compose} ps")


@task(
    help={
        "extra": "Whether extra tooling is included in the infrastructure.",
        "detach": "Whether the infrastructure starts detached.",
    }
)
def start(context, extra=False, detach=False):
    """Start infrastructure."""

    context.run(
        f"{Tools.Compose} {'--profile extra' if extra else ''} up --menu=false --build {'--detach' if detach else ''}"
    )


@task()
def clean(context):
    """Remove all containers, volumes and networks."""

    context.run(f"{Tools.Compose} --profile extra down --volumes")


@task()
def prune(context):
    """Remove all containers, volumes, networks and images."""

    context.run(f"{Tools.Docker} system prune -a --volumes")


@task(variadic=True)
def logs(context, containers):
    """Show logs of all or selected containers."""

    context.run(f"{Tools.Compose} logs -f {containers}")


@task()
def stats(context):
    """Show stats of all containers."""

    context.run(f"{Tools.Compose} stats")
