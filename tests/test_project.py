from grpclib.client import Channel
from testsuite.utils import matching

from api.proto.v1 import admin_pb2


async def test_create_get(project_api, message_to_dict, authorization):
    response = message_to_dict(
        await project_api.CreateProject(
            admin_pb2.CreateProjectRequest(name="Testsuite"),
            metadata=authorization,
        )
    )

    assert response == {
        "project": {
            "id": matching.any_string,
            "name": "Testsuite",
        }
    }

    project_id = response["project"]["id"]

    response = message_to_dict(
        await project_api.GetProjects(
            admin_pb2.GetProjectsRequest(project_ids=[project_id]),
            metadata=authorization,
        )
    )
    assert response == {
        "projects": [
            {
                "id": project_id,
                "name": "Testsuite",
            }
        ]
    }
