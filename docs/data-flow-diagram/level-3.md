```                                            
▄█████ ▄▄ ▄▄  ▄▄▄▄ ▄▄▄▄▄▄ ▄▄▄▄▄ ▄▄   ▄▄    ▄████  ▄▄ ▄▄ ▄▄ ▄▄▄▄  ▄▄▄▄▄ ▄▄    ▄▄ ▄▄  ▄▄ ▄▄▄▄▄
▀▀▀▄▄▄ ▀███▀ ███▄▄   ██   ██▄▄  ██▀▄▀██   ██  ▄▄▄ ██ ██ ██ ██▀██ ██▄▄  ██    ██ ███▄██ ██▄▄ 
█████▀   █   ▄▄██▀   ██   ██▄▄▄ ██   ██    ▀███▀  ▀███▀ ██ ████▀ ██▄▄▄ ██▄▄▄ ██ ██ ▀██ ██▄▄▄ - DFD Version

This section serves as a guideline for breaking down each process in the system into smaller,
    more manageable steps.
It helps give a clearer picture of what happens inside each function and can be used as a
    reference during development to understand the flow and responsibilities of each part.
```

## 1.0 Security Gateway (WAF Layer)

### 1.1 Filter Request
- extractRequestData()
- validateHeaders()
- sanitizeInput()
- forwardRequest()

### 1.2 Rate Limiting
- identifyClient()
- checkRequestCount()
- compareThreshold()
- allowOrBlockRequest()

### 1.3 Threat Logging
- captureThreatData()
- formatLogEntry()
- storeLog()
- notifySystem()


## 2.0 E-Commerce Operations

### 2.1 Authenticate User
- receiveCredentials()
- validateInput()
- verifyUserExists()
- comparePassword()
- generateToken()
- returnAuthResponse()

### 2.2 Manage Products
- receiveProductRequest()
- validateProductData()
- checkProductExists()
- createOrUpdateProduct()
- fetchProductList()
- returnProductResponse()

### 2.3 Process Order
- receiveOrderRequest()
- validateUser()
- validateCartItems()
- checkStockAvailability()
- calculateOrderTotal()
- createOrderRecord()
- updateOrderStatus()
- returnOrderResponse()

### 2.4 Process Payment
- receivePaymentRequest()
- validatePaymentData()
- createPaymentPayload()
- sendToPaymentGateway()
- receivePaymentStatus()
- updatePaymentStatus()
- returnPaymentResult()

### 2.5 Handle Delivery
- receiveDeliveryRequest()
- validateShippingInfo()
- createShipmentOrder()
- sendToCourierService()
- receiveDeliveryStatus()
- updateDeliveryStatus()
- notifyUser()


## 3.0 AI & Analytics

### 3.1 Aggregate Data
- collectOrderData()
- collectUserData()
- collectLogData()
- mergeDatasets()
- cleanData()
- storeAggregatedData()

### 3.2 Detect Anomalies
- loadTransactionData()
- preprocessData()
- applyDetectionModel()
- evaluateAnomalies()
- flagSuspiciousActivity()
- storeDetectionResults()

### 3.3 Forecast Sales
- loadHistoricalSales()
- preprocessTimeSeries()
- applyForecastModel()
- generateForecast()
- evaluatePrediction()
- storeForecastResults()

### 3.4 Generate Reports
- receiveReportRequest()
- fetchAggregatedData()
- formatReportData()
- generateVisualization()
- exportReport()
- deliverReportToAdmin()