import pytest
from google.protobuf import json_format

from api.proto.v1 import admin_grpc, keys_grpc, keyset_grpc

pytest_plugins = ["tests.plugins.plugin"]


@pytest.fixture(name="project_api")
def _project_api_fixture(grpc_channel):
    return admin_grpc.ProjectAPIStub(grpc_channel)


@pytest.fixture(name="message_to_dict")
def _message_to_dict_fixture():
    def convert(message):
        return json_format.MessageToDict(
            message, preserving_proto_field_name=True
        )

    return convert


@pytest.fixture(name="keyset_api")
def _keyset_api_fixture(grpc_channel):
    return keyset_grpc.KeysetAPIStub(grpc_channel)


@pytest.fixture(name="authorization")
def _authrization_fixture():
    return {"authorization": "bearer testsuite-token"}


@pytest.fixture(name="keys_api")
def _keys_api_fixture(grpc_channel):
    return keys_grpc.KeysAPIStub(grpc_channel)
