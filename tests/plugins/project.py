import os
import pathlib
import tempfile

import pytest
import yaml
from testsuite.databases.pgsql import discover

ROOT_DIR_ENV = "ROOT_DIR"
HTTP_SERVER_HOSTPORT = "localhost:18080"
GRPC_SERVER_HOSTPORT = "localhost:18081"


@pytest.fixture(scope="session")
def pgsql_local(project_root, pgsql_local_create):
    discovered = discover.find_schemas(
        "late",
        schema_dirs=[project_root / "postgresql"],
    )
    return pgsql_local_create(list(discovered.values()))


def _yaml_write(path: pathlib.Path, data: dict):
    with path.open("w") as file:
        file.write(
            yaml.dump(
                data,
                Dumper=yaml.CDumper,
                width=80,
                indent=4,
                sort_keys=False,
                allow_unicode=True,
                default_flow_style=False,
            )
        )


def _create_uri(db) -> str:
    if not db.password:
        return f"postgresql://{db.user}@{db.host}:{db.port}/{db.dbname}?sslmode=disable"
    return f"postgresql://{db.user}:{db.passowrd}@{db.host}:{db.port}/{db.dbname}?sslmode=disable"


@pytest.fixture(scope="session")
def project_root():
    if ROOT_DIR_ENV not in os.environ:
        raise Exception("Cannot find ROOT_DIR environment variable")
    return pathlib.Path(os.environ[ROOT_DIR_ENV])


@pytest.fixture(scope="session")
def config_path(project_root, pgsql_local):
    root_config: dict = {}
    root_config_path = project_root / "config.test.yaml"
    if root_config_path.is_file():
        with root_config_path.open() as f:
            root_config = yaml.load(f, Loader=yaml.CLoader)

    root_config.update(
        {
            "postgres": {"uri": _create_uri(pgsql_local["late"])},
            "grpc": {
                "listen": GRPC_SERVER_HOSTPORT,
            },
            "http": {
                "listen": HTTP_SERVER_HOSTPORT,
            },
        }
    )

    config_path = pathlib.Path(tempfile.gettempdir(), "config.yaml")
    _yaml_write(config_path, root_config)

    return config_path


@pytest.fixture(scope="session")
async def service_scope(create_daemon_scope, project_root, config_path):
    async with create_daemon_scope(
        args=[
            f"{project_root}/build/development/late",
            "serve",
            "--config",
            str(config_path),
        ],
        ping_url=f"http://{HTTP_SERVER_HOSTPORT}/v1/health/check",
    ) as scope:
        yield scope


@pytest.fixture
async def service_daemon(
    ensure_daemon_started,
    service_scope,
    mockserver,
    pgsql,
):
    await ensure_daemon_started(service_scope)


@pytest.fixture
def late(create_service_client, service_daemon):
    """
    Returns service client
    """
    return create_service_client(f"http://{HTTP_SERVER_HOSTPORT}")
