import pytest

from pricing import discounted


def test_discount() -> None:
    assert discounted(1000, 15) == 850


def test_rejects_bad_percent() -> None:
    with pytest.raises(ValueError):
        discounted(1000, 120)
