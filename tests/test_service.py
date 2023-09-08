from testsuite.utils import matching


async def test_health(late):
    response = await late.get("/v1/health/check")
    assert response.status == 200

    assert response.json() == {
        "version": matching.any_string,
        "commit": matching.any_string,
        "date": matching.any_string,
    }
