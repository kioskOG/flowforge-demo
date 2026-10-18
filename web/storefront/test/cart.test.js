import { test } from "node:test";
import assert from "node:assert/strict";
import { cartTotal } from "../src/cart.js";

test("adds line items", () => {
  assert.equal(cartTotal([{ qty: 2, price: 150 }, { qty: 1, price: 99 }]), 399);
});

test("empty cart is zero", () => {
  assert.equal(cartTotal([]), 0);
});
