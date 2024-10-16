import ssl

import superinvoke
from superinvoke import Collection
from superinvoke.constants import Platforms

from . import backend, engine, frontend, infra
from .envs import Envs
from .tools import Tools

# Temporal fix very annoying error: certificate verify failed: unable to get local issuer certificate
ssl._create_default_https_context = ssl._create_unverified_context

root = superinvoke.init(tools=Tools, envs=Envs)

root.configure(
    {
        "run": {
            "pty": (Platforms.CURRENT != Platforms.WINDOWS and Envs.Current != Envs.Ci),
        },
    }
)


root.add_collection(
    Collection(
        infra.list,
        infra.start,
        infra.clean,
        infra.prune,
        infra.logs,
        infra.stats,
    ),
    name="infra",
)

root.add_collection(
    Collection(
        backend.lint,
        backend.format,
        backend.test,
        backend.migrate,
        backend.run,
        backend.build,
        backend.deploy,
    ),
    name="backend",
)

root.add_collection(
    Collection(
        engine.lint,
        engine.format,
        engine.test,
        engine.build,
        engine.deploy,
    ),
    name="engine",
)

root.add_collection(
    Collection(
        frontend.lint,
        frontend.format,
        frontend.test,
        frontend.build,
        frontend.deploy,
    ),
    name="frontend",
)
