#!/bin/bash

# ============================================================================
# Product Management API - Manual Testing Commands
# ============================================================================
# 
# This script contains individual curl commands for manually testing the 
# Product Management API endpoints. Each command is independent and can be
# executed separately.
#
# Usage:
#  1. Make the script executable: chmod +x product-api-manual-test.sh
#  2. Run a specific command by uncommenting it and running the script
#  3. Replace placeholders like {product_id} with actual values
#
# Requirements:
#  - curl must be installed
#  - jq is recommended for formatting JSON output (pipe output to jq '.')
#
# ============================================================================

# Set the base URL for the API (modify if needed)
API_URL="http://localhost:8080"

# Uncomment the command you want to run:

# ============================================================================
# 1. Create a new product
# ============================================================================
# curl -X POST "${API_URL}/api/v1/products" \
#   -H "Content-Type: application/json" \
#   -d '{
#     "name": "Biggy CT-250 1 DUS",
#     "cost_price": 31500,
#     "gross_profit_percentage": 10,
#     "shopee_category": "A",
#     "shopee_name": "Manual Shopee Name"
#   }'

# ============================================================================
# 2. Get all products
# ============================================================================
# curl "${API_URL}/api/v1/products"

# ============================================================================
# 3. Get all products with pagination
# ============================================================================
# curl "${API_URL}/api/v1/products?page_size=10&page_number=0"

# ============================================================================
# 4. Get products with name filter
# ============================================================================
# curl "${API_URL}/api/v1/products?name=Test"

# ============================================================================
# 5. Get a specific product by ID
# ============================================================================
# Replace {product_id} with an actual UUID from your database
# curl "${API_URL}/api/v1/products/64e9df73-f92f-49e0-a393-135e480e89ed"

# ============================================================================
# 6. Update a product
# ============================================================================
# Replace {product_id} with an actual UUID from your database
# curl -X PUT "${API_URL}/api/v1/products/{product_id}" \
#   -H "Content-Type: application/json" \
#   -d '{
#     "name": "Updated Product Name",
#     "gross_profit_percentage": 30.75,
#     "shopee_category": "B",
#     "shopee_name": "Manual Shopee Name Updated"
#   }'

# ============================================================================
# 7. Delete a product (soft delete)
# ============================================================================
# Replace {product_id} with an actual UUID from your database
# curl -X DELETE "${API_URL}/api/v1/products/{product_id}"

# ============================================================================
# 8. Test invalid create request (validation should fail)
# ============================================================================
# curl -X POST "${API_URL}/api/v1/products" \
#   -H "Content-Type: application/json" \
#   -d '{
#     "name": "Invalid Product",
#     "cost_price": 100,
#     "gross_profit_percentage": 200,
#     "shopee_category": "Z",
#     "shopee_name": ""
#   }'

# ============================================================================
# 9. Test invalid update request (validation should fail)
# ============================================================================
# Replace {product_id} with an actual UUID from your database
# curl -X PUT "${API_URL}/api/v1/products/{product_id}" \
#   -H "Content-Type: application/json" \
#   -d '{
#     "name": "",
#     "gross_profit_percentage": -10,
#     "shopee_category": "X"
#   }'

# ============================================================================
# 10. Test getting a non-existent product
# ============================================================================
# curl "${API_URL}/api/v1/products/00000000-0000-0000-0000-000000000000"

# ============================================================================
# Helper functions for common testing patterns
# ============================================================================

# Create a product and return the ID
create_and_get_id() {
  response=$(curl -s -X POST "${API_URL}/api/v1/products" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Test Product For ID",
      "cost_price": 100.50,
      "gross_profit_percentage": 25.50,
      "shopee_category": "A",
      "shopee_name": "Test Shopee Name For ID"
    }')
  
  # Get all products and extract the first ID
  all_products=$(curl -s "${API_URL}/api/v1/products")
  product_id=$(echo $all_products | jq -r '.data[0].id')
  echo $product_id
}

# Uncomment to create a product and get its ID for testing
# ID=$(create_and_get_id)
# echo "Created product with ID: $ID"
# 
# # Now you can use the ID in subsequent commands
# curl "${API_URL}/api/v1/products/$ID"