import pytest
from google.protobuf import json_format

from api.proto.v1 import admin_grpc

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
