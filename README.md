# [Capitol News Illinios](https://capitolnewsillinois.com/) web scraper. 

This project will scrape [Capitol News Illinios](https://capitolnewsillinois.com/) looking for content that has been posted today. After content has been found, it will send the scraped content to an endpoint of your choosing. 

### Running the project

``` git clone https://github.com/zepez/capitol-news.git ```

``` cd capitol-news ```

Change the environment variables to suit your needs: 
  - TZ : Your timezone. Defalts to America/New_York
  - cron : Any valid cron interval. Recommended to leave as once per day. The script will find anything that was posted within the date that it is running. 
  - seconds : How long do you want to look for new content? The longer this value is, the more requests will be made. 
  - endpoint : This is the API endpoint that the content will be sent.

  ``` docker build . ```

  ``` docker run CONTAINERID ```

