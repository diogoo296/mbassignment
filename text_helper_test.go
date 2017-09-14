// Tests TextHelperInit() subfunctions
package main

import (
  "fmt"
  "strings"
  "testing"
)

var shortPlainMsg string = strings.Repeat("a", 10)
var shortSpecialMsg string = "\f\\~[]{}|€"
var maxShortPlainMsg string = strings.Repeat("a", 160)
var longPlainMsg string = strings.Repeat("a", 170)
var veryLongPlainMsg string = strings.Repeat("a", 612)
var wayTooLongPlainMsg string = strings.Repeat("a", 39016)
// 169 normal chars + 1 (last) special char
var longPlainMsgEndSpecial string =
  fmt.Sprintf("%s€", strings.Repeat("a", 169))
var shortUnicodeMsg string = "日本語€g"
var maxShortUnicodeMsg string = strings.Repeat("日", 70)
var longUnicodeMsg string = strings.Repeat("日", 80)
var veryLongUnicodeMsg string = strings.Repeat("日", 268)
var wayTooLongUnicodeMsg string = strings.Repeat("日", 17086)

func checkPlainText(msg string, expected bool, t *testing.T) {
  tHelper, err := TextHelperInit(msg);
  if err != nil {
    t.Errorf("Expected error to be nil: %#v", err)
  } else if tHelper.PlainText != expected {
    t.Errorf(
      "Expected TextHelper.PlanText=%t. Full object: %#v",
      expected, tHelper)
  }
}

func TestSetPlainText(t *testing.T) {
  // Test simple text
  checkPlainText(shortPlainMsg, true, t)
  // Test special characters
  checkPlainText(shortSpecialMsg, true, t)
  // Test unicode message
  checkPlainText(shortUnicodeMsg, false, t)
}

func checkNumChars(msg string, expected int, t *testing.T) {
  tHelper, err := TextHelperInit(msg);
  if err != nil {
    t.Errorf("Expected error to be nil: %#v", err)
  } else if tHelper.NumChars != expected {
    t.Errorf(
      "Expected TextHelper.NumChars=%d. Full object: %#v",
      expected, tHelper)
  }
}

func TestCountChars(t *testing.T) {
  // Test simple text
  checkNumChars(shortPlainMsg, 10, t)
  // Test special characters
  checkNumChars(shortSpecialMsg, 18, t)
  // Test unicode message
  checkNumChars(shortUnicodeMsg, 5, t)
}

func checkPartSize(msg string, expected int, t *testing.T) {
  tHelper, err := TextHelperInit(msg);
  if err != nil {
    t.Errorf("Expected error to be nil: %#v", err)
  } else if tHelper.PartSize != expected {
    t.Errorf(
      "Expected TextHelper.PartSize=%d. Full object: %#v",
      expected, tHelper)
  }
}

func TestSetPartSize(t *testing.T) {
  // Plain text messages
  checkPartSize(shortPlainMsg, PLAIN_SMS_MAX_LEN, t)
  checkPartSize(shortSpecialMsg, PLAIN_SMS_MAX_LEN, t)
  checkPartSize(maxShortPlainMsg, PLAIN_SMS_MAX_LEN, t)
  checkPartSize(longPlainMsg, PLAIN_CSMS_MAX_LEN, t)
  checkPartSize(longPlainMsgEndSpecial, PLAIN_CSMS_MAX_LEN, t)

  // Unicode messages
  checkPartSize(shortUnicodeMsg, UNICODE_SMS_MAX_LEN, t)
  checkPartSize(maxShortUnicodeMsg, UNICODE_SMS_MAX_LEN, t)
  checkPartSize(longUnicodeMsg, UNICODE_CSMS_MAX_LEN, t)
}

func checkNumParts(msg string, expected int, t *testing.T) {
  tHelper, err := TextHelperInit(msg);
  if err != nil {
    t.Errorf("Expected error to be nil: %#v", err)
  } else if tHelper.NumParts != expected {
    t.Errorf(
      "Expected TextHelper.NumParts=%d. Full object: %#v",
      expected, tHelper)
  }
}

func expectError(msg string, t *testing.T) {
  if _, err := TextHelperInit(msg); err == nil {
    t.Error("Expected error to not be nil.")
  }
}

func TestSplitBody(t *testing.T) {
  // Plain text messages
  checkNumParts(shortPlainMsg, 1, t)
  checkNumParts(shortSpecialMsg, 1, t)
  checkNumParts(maxShortPlainMsg, 1, t)
  checkNumParts(longPlainMsg, 2, t)
  checkNumParts(longPlainMsgEndSpecial, 2, t)
  checkNumParts(veryLongPlainMsg, 4, t)
  expectError(wayTooLongPlainMsg, t)

  // Unicode messages
  checkNumParts(shortUnicodeMsg, 1, t)
  checkNumParts(maxShortUnicodeMsg, 1, t)
  checkNumParts(longUnicodeMsg, 2, t)
  checkNumParts(veryLongUnicodeMsg, 4, t)
  expectError(wayTooLongUnicodeMsg, t)
}
