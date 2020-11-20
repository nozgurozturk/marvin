/**
 *
 * @param {string} value // input
 * @param {number} charLength // length of input
 * @param {field} field  // input name
 */

export const requireMinCharacter = (value: string, charLength: number, field: string) => {
  if (value.length >= charLength) {
    return true;
  }

  return `${field} need to have at least ${charLength} character long`;
};
