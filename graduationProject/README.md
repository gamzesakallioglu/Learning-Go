Basket service developed with Golang.     
You can use the endpoints above to use this application.

**Sign-up & Sign-in**


**Sign Up as a customer**         
>localhost:8080/api/v1/sign-up      

     Request Body:       
     {            
         "name": "some name",         
         "email": "some@gmail.com",             
         "password": "somepassword123",
         "phoneNumber": "555 555 55 55",
         "address": "some address"
     }
     
          
               
**Sign-in as a customer**
>localhost:8080/api/v1/login     

    Request Body:
    {
        "email": "some@gmail.com",
        "password": "somepassword123"
    }

This request will give a token as a reponse. You can use that token to access endpoints accessible for customers by adding Authorization-Token (key-value) to request header.                      
            
           

**Sign-up as a user**                             
Only admin users can create new users.
>localhost:8080/api/v1/user/sign-up                      
                    
         Request Body:        
         {          
              "name": "Some Name",                                
              "email": "some@gmail.com",
              "password": "123456",
              "phoneNumber": "555 555 55 55"
          }
          
          
**Sign-in as a user**
>localhost:8080/api/v1/user/login         
               
         Request Body:        
         {          
              "email": "some@gmail.com",
              "password": "123456"
          }

Passwords are stored in table with Md5 hashing.

          
               
**Product and Product Category Create**                
Only admin users can create product and product category.                       

In order to create product categories, an admin user should make a POST request to 
>localhost:8080/api/v1/productCategories  (POST)              
     
     
Request body should be a form with a csv file (key: file)        
This file will be read line by line and will insert the line into database if there isn’t any product category with same name with the line, if there is; the record will be updated.         
     
     

In order to create product, an admin user should make a POST request to         
>localhost:8080/api/v1/products  (POST)             

          Request Body:
          {
              "name": "new product",
              "description": "some desc",
              "price": 710.65,
              "stockCode": "NW005AJ54",
              "stockNumber": 700,
              "category": {
                  "id": "4be2f5e7-833b-4f46-8ce9-02a3f39bc74b",
                  "name": "Health & Wellness products"
              }
           }

                    
**Product and Category Listing**       
There isn’t any authorization middleware to use this endpoint. Both users, customers and guests can        
list products and product categories. Both endpoints use pagination for better performance.             
>localhost:8080/api/v1/products (GET)    

     {        
        "page": 1,       
        "pageSize": 10,     
        "sorting": "ID desc",
        "total_rows": 6,
        "total_pages": 1,
        "rows": [
            {
                "id": "2c5c2c18-de40-459a-bd19-787f46140e43",
                "name": "Electronics",
                "parentID": "0"
            },......           
       
          
>localhost:8080/api/v1/productCategories (GET)   

     {       
         "page": 1,         
         "pageSize": 10,           
         "sorting": "ID desc",
         "total_rows": 3,
         "total_pages": 1,
         "rows": [
             {
                 "category": {
                     "description": "Electronics products",
                     "id": "2c5c2c18-de40-459a-bd19-787f46140e43",
                     "isParent": true,
                     "name": "Electronics",
                     "parentID": "0"
                 },.......

                              
**Delete Product & Product Category**      
Of course, only admin users can delete a product or a product category.     
>localhost:8080/api/v1/products/:id       
>localhost:8080/api/v1/productCategories/:id     
                         
                         
**Shopping Cart**
Only customers can access and manage their own shopping cart. They can add an item, increase or decrease the amount of the item, list the items or delete an item.
                              
                              
**Add An Item**                   
>localhost:8080/api/v1/cart  (POST) 


     {      
        "product":{      
             "id": "2ff597b3-24db-4714-ab77-70fb12ccd496",    
             "name": "Dell PC"   
         },    
         "quantity": 2    
     }     
          
          
**Delete An Item**   
>localhost:8080/api/v1/cart/item/:id  (DELETE)             
                    
                    
**List The Items**     
>localhost:8080/api/v1/cart   (GET)     
                    
                    
**Order**
Customers can complete the order with what they have in their shopping cart. They can cancel the order within 14 days their order date or list their past orders.
                              
                              
**Complete Order**  
>localhost:8080/api/v1/orders   (POST)     

    {        
        "phoneNumber": "555 555 55 55",        
        "address": "some address"         
    }         

    Response: "Order has been completed. Order number: 8c3f01e6071e94d73f3709f35370112d"     
                    
                    
**Cancel Order**         
>localhost:8080/api/v1/orders/cancel/:id   (GET)      
                              
                              
**List Past Orders**           
>localhost:8080/api/v1/orders   (GET)           

    [
     {
         "address": "some address",
         "cancelDate": "0001-01-01 03:00:00 +0300 +03",
         "id": "775444e3-848c-4f90-88b9-3820e64132e0",
         "orderDate": "2022-04-18 06:48:57",
         "orderStatus": "awaiting shipment",
         "orderTotal": 22005.5,
         "paymentDate": "2022-04-18 06:48:57",
         "paymentDueDate": "2022-04-19 06:48:57",
         "paymentStatus": "completed",
         "phoneNumber": "555 555 55 55",
         "products": [
             {
                 "pricePerItem": 2000.5,
                 "product": {
                     "category": {
                         "description": "Electronics products",
                         "id": "2c5c2c18-de40-459a-bd19-787f46140e43",
                         "isParent": true,
                         "name": "Electronics",
                         "parentID": "0"
                     },
                     "description": "Dell PC 15.6 inches",
                     "id": "bccd451c-56f5-4863-ae4a-4b921fda906e",
                     "name": "Dell PC",
                     "price": 2000.5,
                     "stockCode": "TECHA1B50001",
                     "stockNumber": 4920
                 },
                 "quantity": 4
             },
             {
                 "pricePerItem": 2000.5,
                 "product": {
                     "category": {
                         "description": "Electronics products",
                         "id": "2c5c2c18-de40-459a-bd19-787f46140e43",
                         "isParent": true,
                         "name": "Electronics",
                         "parentID": "0"
                     },
                     "description": "Dell PC 15.6 inches",
                     "id": "2ff597b3-24db-4714-ab77-70fb12ccd496",
                     "name": "Dell PC",
                     "price": 2000.5,
                     "stockCode": "TECHA1B50001",
                     "stockNumber": 4917
                 },
                 "quantity": 7
             }
         ],
         "receiveDate": "0001-01-01 03:00:00",
         "shippingDate": "0001-01-01 03:00:00"
     },.......
     ]
