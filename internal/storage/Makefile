

upgrade-from-schema-sqlite3:
	# remove old database
	rm -f /tmp/data.db
	# populate the database
	go run generate/main.go -file=/tmp/data.db
	# create boilerplate code
	sqlboiler sqlite3 --no-rows-affected  --no-context