## Build
  - `docker-compose up` from the root directory will build and run the project.
## Usage
  - Send a post request to "localhost:8000/getProductInfo" to scrape and store data. Format of the body is described below.
  - Send a get request to "localhost:8010/getAllProducts" to view all the stored data.
  - NOTE : On some products ([Like this one](https://www.amazon.in/dp/B07J1ZBCZJ/ref=s9_acsd_ri_bw_c2_x_0_i?pf_rd_m=A1K21FY43GMZF8&pf_rd_s=merchandised-search-5&pf_rd_r=PSSGG284GH10YC2MSW4Y&pf_rd_t=101&pf_rd_p=33ff9bf4-e55c-4ab2-8f44-1458082fb86c&pf_rd_i=17941593031)) the service may not be able to scrape all the necessary information as they have different document structure. The rest of the product pages ([Like this one](https://www.amazon.in/dp/B08R6NFZ6R?ref_=nav_em_nav-pc-ftvlite_0_2_3_2)) work fine.
## Service 1
  - Hosted at localhost:8000
  - Endpoints
    - `/getProductInfo`
      - Method : POST
      - Request body format
        ```
        {
            "url": "page_link_here"
        }
      - Response
        - 200 : Product information is scrapped, stored using the second service and returned in the response
      - Errors
        - 400 : Invalid domain (Allowed : "amazon.in", "www.amazon.in")
        - 500 : Other errors
## Service 2
  - Hosted at localhost:8010
  - Endpoints
    - `/writeProductInfo`
      - Method : POST
      - Used by Service 1
    - `/getAllProducts`
      - Method : GET
      - Response
        - 200 : An array of JSONs describing all the products.
      - Error
        - 500 : Server error