import pytest
from grpclib.client import Channel
from testsuite.utils import matching

from api.proto.v1 import admin_pb2, keys_pb2, keyset_pb2


@pytest.fixture(name="create_project")
def _create_project(project_api, message_to_dict, authorization):
    async def call(name: str) -> str:
        project_response = message_to_dict(
            await project_api.CreateProject(
                admin_pb2.CreateProjectRequest(name="Testsuite"),
                metadata=authorization,
            )
        )

        return project_response["project"]["id"]

    return call


@pytest.fixture(name="create_keyset")
def _create_keyset(keyset_api, message_to_dict, authorization):
    async def call(project_id: str, name: str) -> str:
        keyset_response = message_to_dict(
            await keyset_api.CreateKeyset(
                keyset_pb2.CreateKeysetRequest(
                    project_id=project_id,
                    name=name,
                ),
                metadata=authorization,
            )
        )

        return keyset_response["keyset"]["id"]

    return call


async def test_create_get_key(
    create_project, create_keyset, keys_api, message_to_dict, authorization
):
    KEY = "test_key"

    project_id = await create_project("test")
    keyset_id = await create_keyset(project_id, "test")

    key_response = message_to_dict(
        await keys_api.CreateKey(
            keys_pb2.CreateKeyRequest(
                key=KEY,
                keyset_id=keyset_id,
                translations=[
                    keys_pb2.Translation(
                        language="en",
                        texts=["apple", "apples"],
                    ),
                ],
            ),
            metadata=authorization,
        ),
    )

    get_keys_response = message_to_dict(
        await keys_api.GetKeys(
            keys_pb2.GetKeysRequest(keyset_id=keyset_id),
            metadata=authorization,
        )
    )

    assert get_keys_response == {
        "keys": [
            {
                "keyset_id": keyset_id,
                "key": KEY,
                "translations": [
                    {
                        "language": "en",
                        "texts": ["apple", "apples"],
                    }
                ],
            }
        ]
    }


async def test_create_delete_get(
    create_project, create_keyset, keys_api, message_to_dict, authorization
):
    KEY = "test_key"

    project_id = await create_project("test")
    keyset_id = await create_keyset(project_id, "test")

    key_response = message_to_dict(
        await keys_api.CreateKey(
            keys_pb2.CreateKeyRequest(
                key=KEY,
                keyset_id=keyset_id,
                translations=[
                    keys_pb2.Translation(
                        language="en",
                        texts=["apple", "apples"],
                    ),
                ],
            ),
            metadata=authorization,
        ),
    )

    delete_response = message_to_dict(
        await keys_api.DeleteKey(
            keys_pb2.DeleteKeyRequest(keyset_id=keyset_id, key=KEY),
            metadata=authorization,
        )
    )

    get_keys_response = message_to_dict(
        await keys_api.GetKeys(
            keys_pb2.GetKeysRequest(keyset_id=keyset_id),
            metadata=authorization,
        )
    )

    assert get_keys_response == {}


async def test_create_update(
    create_project, create_keyset, keys_api, message_to_dict, authorization
):
    KEY = "test_key"

    project_id = await create_project("test")
    keyset_id = await create_keyset(project_id, "test")

    key_response = message_to_dict(
        await keys_api.CreateKey(
            keys_pb2.CreateKeyRequest(
                key=KEY,
                keyset_id=keyset_id,
                translations=[
                    keys_pb2.Translation(
                        language="en",
                        texts=["apple", "apples"],
                    ),
                ],
            ),
            metadata=authorization,
        ),
    )

    update_response = message_to_dict(
        await keys_api.UpdateKey(
            keys_pb2.UpdateKeyRequest(
                keyset_id=keyset_id,
                key=KEY,
                translations=[
                    keys_pb2.Translation(
                        language="en",
                        texts=["apple", "apples"],
                    ),
                    keys_pb2.Translation(
                        language="ru",
                        texts=["яблоко", "яблока", "яблок"],
                    ),
                ],
            ),
            metadata=authorization,
        )
    )

    get_keys_response = message_to_dict(
        await keys_api.GetKeys(
            keys_pb2.GetKeysRequest(keyset_id=keyset_id),
            metadata=authorization,
        )
    )

    assert get_keys_response == {
        "keys": [
            {
                "keyset_id": keyset_id,
                "key": KEY,
                "translations": [
                    {"language": "en", "texts": ["apple", "apples"]},
                    {"language": "ru", "texts": ["яблоко", "яблока", "яблок"]},
                ],
            }
        ]
    }
