#!/bin/bash

# Base URL
BASE_URL="http://localhost:8080"

# Set colors for better output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}===== Product Management API Testing =====\n${NC}"

# 1. Create a product
echo -e "${BLUE}Testing: Create a product (POST /api/v1/products)${NC}"
CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product A",
    "cost_price": 100.50,
    "gross_profit_percentage": 25.50,
    "shopee_category": "A"
  }')

if [ -z "$CREATE_RESPONSE" ]; then
  echo -e "${GREEN}Product created successfully!${NC}"
else
  echo -e "${RED}Creation failed with response:${NC}"
  echo $CREATE_RESPONSE | jq '.'
fi

# 2. Create a second product for testing
echo -e "\n${BLUE}Creating another product for testing...${NC}"
PRODUCT_B_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product B",
    "cost_price": 200.75,
    "gross_profit_percentage": 30.00,
    "shopee_category": "B"
  }')

# 3. Get all products
echo -e "\n${BLUE}Testing: Get all products (GET /api/v1/products)${NC}"
GET_ALL_RESPONSE=$(curl -s "${BASE_URL}/api/v1/products")
echo $GET_ALL_RESPONSE | jq '.'

# Extract the first product ID for further operations
PRODUCT_ID=$(echo $GET_ALL_RESPONSE | jq -r '.data[0].id')
echo -e "Using product ID: ${PRODUCT_ID}"

if [ "$PRODUCT_ID" = "null" ]; then
  echo -e "${RED}Failed to get a valid product ID. Exiting tests.${NC}"
  exit 1
fi

# 4. Get a specific product
echo -e "\n${BLUE}Testing: Get product by ID (GET /api/v1/products/$PRODUCT_ID)${NC}"
GET_PRODUCT_RESPONSE=$(curl -s "${BASE_URL}/api/v1/products/$PRODUCT_ID")
echo $GET_PRODUCT_RESPONSE | jq '.'

# 5. Update a product
echo -e "\n${BLUE}Testing: Update product (PUT /api/v1/products/$PRODUCT_ID)${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "${BASE_URL}/api/v1/products/$PRODUCT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product Name",
    "gross_profit_percentage": 35.75,
    "shopee_category": "C"
  }')

if [ -z "$UPDATE_RESPONSE" ]; then
  echo -e "${GREEN}Product updated successfully!${NC}"
else
  echo -e "${RED}Update failed with response:${NC}"
  echo $UPDATE_RESPONSE | jq '.'
fi

# 6. Verify the update
echo -e "\n${BLUE}Verifying product update...${NC}"
VERIFY_UPDATE=$(curl -s "${BASE_URL}/api/v1/products/$PRODUCT_ID")
echo $VERIFY_UPDATE | jq '.'

# 7. Get products with pagination and filtering
echo -e "\n${BLUE}Testing: Get products with pagination and filtering (GET /api/v1/products?page_size=10&page_number=0&name=Updated)${NC}"
FILTERED_RESPONSE=$(curl -s "${BASE_URL}/api/v1/products?page_size=10&page_number=0&name=Updated")
echo $FILTERED_RESPONSE | jq '.'

# 8. Delete a product
echo -e "\n${BLUE}Testing: Delete product (DELETE /api/v1/products/$PRODUCT_ID)${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/api/v1/products/$PRODUCT_ID")

if [ -z "$DELETE_RESPONSE" ]; then
  echo -e "${GREEN}Product deleted successfully!${NC}"
else
  echo -e "${RED}Deletion failed with response:${NC}"
  echo $DELETE_RESPONSE | jq '.'
fi

# 9. Verify the delete by trying to fetch the deleted product
echo -e "\n${BLUE}Verifying product deletion...${NC}"
VERIFY_DELETE=$(curl -s "${BASE_URL}/api/v1/products/$PRODUCT_ID")
echo $VERIFY_DELETE | jq '.'

# 10. Test with invalid data (validation errors)
echo -e "\n${BLUE}Testing: Create product with invalid data (validation errors)${NC}"
INVALID_CREATE=$(curl -s -X POST "${BASE_URL}/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "A very long product name that exceeds the maximum allowed length of fifty characters",
    "cost_price": -100,
    "gross_profit_percentage": "invalid",
    "shopee_category": "Z"
  }')
echo $INVALID_CREATE | jq '.'

echo -e "\n${GREEN}API Testing Completed!${NC}"