package order

// GraphQL queries and mutations for order creation

const CreateOrderMutation = `
mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) {
    order {
      id
      druid
      status
      pricing {
        totalPrice
        basePrice
        discounts {
          type
          amount
        }
      }
    }
    errors {
      field
      message
    }
  }
}
`

const ValidateOrderMutation = `
mutation ValidateOrder($input: ValidateOrderInput!) {
  validateOrder(input: $input) {
    valid
    errors {
      field
      message
    }
    warnings {
      field
      message
    }
  }
}
`

const GetPricingQuery = `
query GetOrderPricing($input: PricingInput!) {
  getOrderPricing(input: $input) {
    basePrice
    totalPrice
    discounts {
      type
      amount
      description
    }
    availableVehicleTypes {
      id
      name
      price
    }
  }
}
`

const GetVehicleTypesQuery = `
query GetVehicleTypes {
  vehicleTypes {
    id
    name
    description
    basePrice
    capabilities
  }
}
`

const GetCapabilitiesQuery = `
query GetCapabilities {
  capabilities {
    id
    name
    description
    price
    category
  }
}
`
