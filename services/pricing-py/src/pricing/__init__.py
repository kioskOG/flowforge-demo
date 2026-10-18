"""Discount rules for the storefront."""


def discounted(price_cents: int, percent: int) -> int:
    """Price after a percentage discount, rounded down to the cent."""
    if not 0 <= percent <= 100:
        raise ValueError("percent must be between 0 and 100")
    return price_cents * (100 - percent) // 100
