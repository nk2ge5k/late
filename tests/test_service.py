from testsuite.utils import matching


async def test_health(late):
    """
    Testing that the /v1/health/check handler returns valid data
    """

    response = await late.get("/healthz")
    assert response.status == 200

    assert response.json() == {
        "version": matching.any_string,
        "commit": matching.any_string,
        "date": matching.any_string,
    }


async def test_metrics(late):
    """
    Testing that metrics are accessible via an HTTP request.
    """

    response = await late.get("/metrics")
    assert response.status == 200
