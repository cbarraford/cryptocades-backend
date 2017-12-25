# lint/syntax-check apib file
validate-api: 
	docker run --rm -v ${CURDIR}:/docs apiaryio/drafter drafter -ul /docs/api.apib

# Create the swagger API specification that is consumed by test/api.
#
# Uses freely available online transformer
# https://apimatic.io/transformer (see documentation at bottom of the page)
build-swagger: validate-api
	curl -sX POST "https://apitransformer.com/api/transform?output=swagger20json" -F upload=@api.apib -o swagger.json

# Create an HTML file containing the API documentation.
# Open out/api.html in your browser.
build-html: validate-api
	docker run --rm -v ${CURDIR}:/docs humangeo/aglio:2.2.0-f -i /docs/api.apib -o /docs/api.html
