package api

import (
	"errors"
	"fmt"
	"strings"
)

// LoginParamsValidation validates login parameters.
func LoginParamsValidation(mobileNumber, userID, pan, mpin, password string) (map[string]string, error) {
	outDict := make(map[string]string)

	if mobileNumber != "" {
		if len(mobileNumber) != 10 && len(mobileNumber) != 13 {
			return nil, errors.New("input Number length must be 13 (with country code (+91)) or 10 (without country code)")
		}
		if len(mobileNumber) == 10 {
			mobileNumber = "+91" + mobileNumber
		}
		outDict["mobileNumber"] = mobileNumber
	} else if pan != "" {
		pan = strings.ToUpper(pan)
		if len(pan) != 10 {
			return nil, errors.New("validation Errors! Length of PAN should be 10")
		}
		if !isValidPAN(pan) {
			return nil, errors.New("validation Errors! Given PAN Number is Not Valid")
		}
		outDict["pan"] = pan
	} else if userID != "" {
		outDict["userID"] = userID
	} else {
		return nil, errors.New("validation Errors! Pass any of Mobile Number, User ID or Pan")
	}

	if mpin != "" {
		outDict["mpin"] = mpin
	} else if password != "" {
		outDict["password"] = password
	}

	return outDict, nil
}

func isValidPAN(pan string) bool {
	if len(pan) != 10 {
		return false
	}
	for i := 0; i < 5; i++ {
		if !isAlpha(pan[i]) {
			return false
		}
	}
	for i := 5; i < 9; i++ {
		if !isDigit(pan[i]) {
			return false
		}
	}
	if !isAlpha(pan[9]) {
		return false
	}
	return true
}

func isAlpha(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// ValidateConfiguration validates the consumer key and secret.
func ValidateConfiguration(consumerKey, consumerSecret string) error {
	if consumerKey == "" {
		return errors.New("please provide the Consumer Key parameter. Without Consumer Key the API cannot be accessed")
	}
	if consumerSecret == "" {
		return errors.New("please provide the Consumer Secret parameter. Without Consumer Secret the API cannot be accessed")
	}
	return nil
}

// PlaceOrderValidation validates order parameters.
func PlaceOrderValidation(exchangeSegment, product, price, orderType, quantity, validity, tradingSymbol, transactionType, amo, disclosedQuantity, marketProtection, pf, triggerPrice, tag string) error {
	if err := validateString(exchangeSegment, exchangeSegmentAllowedValues, "Exchange segment"); err != nil {
		return err
	}
	if err := validateString(product, productAllowedValues, "Product"); err != nil {
		return err
	}
	if err := validateString(price, nil, "Price"); err != nil {
		return err
	}
	if err := validateString(orderType, orderTypeAllowedValues, "Order type"); err != nil {
		return err
	}
	if err := validateString(quantity, nil, "Quantity"); err != nil {
		return err
	}
	if err := validateString(validity, []string{"DAY", "IOC"}, "Validity"); err != nil {
		return err
	}
	if err := validateString(tradingSymbol, nil, "Trading symbol"); err != nil {
		return err
	}
	if err := validateString(transactionType, []string{"B", "S", "Buy", "Sell"}, "Transaction type"); err != nil {
		return err
	}
	if amo != "" {
		if err := validateString(amo, nil, "AMO"); err != nil {
			return err
		}
	}
	if disclosedQuantity != "" {
		if err := validateString(disclosedQuantity, nil, "Disclosed quantity"); err != nil {
			return err
		}
	}
	if marketProtection != "" {
		if err := validateString(marketProtection, nil, "Market protection"); err != nil {
			return err
		}
	}
	if pf != "" {
		if err := validateString(pf, nil, "PF"); err != nil {
			return err
		}
	}
	if triggerPrice != "" {
		if err := validateString(triggerPrice, nil, "Trigger price"); err != nil {
			return err
		}
	}
	if tag != "" {
		if err := validateString(tag, nil, "Tag"); err != nil {
			return err
		}
	}
	return nil
}

func validateString(value string, allowedValues []string, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s must be a string", fieldName)
	}
	if allowedValues != nil && !contains(allowedValues, value) {
		return fmt.Errorf("invalid %s. Allowed values are %v", fieldName, allowedValues)
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// CancelOrderValidation validates the order ID for cancellation.
func CancelOrderValidation(orderID, amo string) error {
	if orderID == "" || strings.TrimSpace(orderID) == "" {
		return errors.New("order_id parameter must be a non-empty string")
	}
	if amo != "" {
		if err := validateString(amo, nil, "AMO"); err != nil {
			return err
		}
	}
	return nil
}

// OrderHistoryValidation validates the order ID for order history.
func OrderHistoryValidation(orderID string) error {
	if orderID == "" || strings.TrimSpace(orderID) == "" {
		return errors.New("order_id parameter must be a non-empty string")
	}
	return nil
}

// MarginValidation validates margin parameters.
func MarginValidation(exchangeSegment, price, orderType, product, quantity, instrumentToken, transactionType, triggerPrice string) error {
	if err := validateString(exchangeSegment, exchangeSegmentAllowedValues, "Exchange segment"); err != nil {
		return err
	}
	if err := validateString(product, productAllowedValues, "Product"); err != nil {
		return err
	}
	if err := validateString(price, nil, "Price"); err != nil {
		return err
	}
	if err := validateString(orderType, orderTypeAllowedValues, "Order type"); err != nil {
		return err
	}
	if err := validateString(quantity, nil, "Quantity"); err != nil {
		return err
	}
	if err := validateString(instrumentToken, nil, "Instrument token"); err != nil {
		return err
	}
	if err := validateString(transactionType, []string{"B", "S", "Buy", "Sell", "sell", "buy"}, "Transaction type"); err != nil {
		return err
	}
	if triggerPrice != "" {
		if err := validateString(triggerPrice, nil, "Trigger price"); err != nil {
			return err
		}
	}
	return nil
}

// LimitsValidation validates limit parameters.
func LimitsValidation(segment, exchange, product string) error {
	if err := validateString(segment, segmentLimits, "Segment"); err != nil {
		return err
	}
	if err := validateString(exchange, exchangeLimits, "Exchange"); err != nil {
		return err
	}
	if err := validateString(product, productLimits, "Product"); err != nil {
		return err
	}
	return nil
}
