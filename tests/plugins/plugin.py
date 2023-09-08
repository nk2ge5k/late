import pytest

pytest_plugins = [
    "testsuite.pytest_plugin",
    "testsuite.databases.pgsql.pytest_plugin",
    "testsuite.plugins.mocked_time",
    "testsuite.plugins.object_hook",
    "tests.plugins.project",
]
