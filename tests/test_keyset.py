from grpclib.client import Channel
from testsuite.utils import matching

from api.proto.v1 import admin_pb2, keyset_pb2


async def test_create_get(project_api, keyset_api, message_to_dict):
    NAME = "Test keyset"

    project_response = message_to_dict(
        await project_api.CreateProject(
            admin_pb2.CreateProjectRequest(name="Testsuite")
        )
    )

    project_id = project_response["project"]["id"]

    response = message_to_dict(
        await keyset_api.CreateKeyset(
            keyset_pb2.CreateKeysetRequest(
                project_id=project_id,
                name=NAME,
            )
        )
    )

    assert "keyset" in response
    keyset = response["keyset"]

    assert keyset == {
        "id": matching.any_string,
        "project_id": project_id,
        "name": NAME,
    }

    response = message_to_dict(
        await keyset_api.GetKeysets(
            keyset_pb2.GetKeysetsRequest(
                project_id=project_id,
            )
        )
    )
    assert "keysets" in response

    keysets = response["keysets"]
    assert len(keysets) == 1
    assert keysets.pop() == keyset
