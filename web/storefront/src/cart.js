/** Sum of line items, in cents. */
export function cartTotal(items) {
  return items.reduce((sum, { qty, price }) => sum + qty * price, 0);
}
