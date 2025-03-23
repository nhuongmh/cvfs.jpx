import tempfile
import traceback
from datetime import datetime
from newspaper import Article
from cvfspy.usecases.common import BaseSchema

# fetch the article 
# Gen JSON file of the article object


class SimpleArticle(BaseSchema):
    title: str
    content: str = ""
    image: str = ""
    author: str = ""
    publish_date: str | datetime | None = ""
    origin: str = ""
    summary: str = ""
    keywords: list[str] = []


def fetch_article(url) -> SimpleArticle:
    article = Article(url)
    article.download()
    article.parse()
    
    return SimpleArticle(
        title=article.title,
        content=article.text,
        image=article.top_image,
        author=", ".join(article.authors),
        publish_date=article.publish_date,
        origin=url,
        summary=article.summary,
        keywords=article.keywords
    )

# if __name__ == '__main__':
#     try:
#         print(f"Fetching article from {args.url}")
#         article = fetch_article(args.url)
#         if article:
#             # save to temporary json file
#             temp_file = tempfile.NamedTemporaryFile(prefix='article-', suffix='.json', delete=False)
#             temp_file.write(article.model_dump_json(indent=2).encode())
#             temp_file.close()
#             print(f"ARTICLE_JSON: {temp_file.name}")
#         else:
#             print('ERROR: Failed to fetch article')
#     except Exception as e:
#         traceback.print_exc()
#         print('ERROR: Failed to fetch article')
#         exit(1)
