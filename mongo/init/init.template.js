// создаем пользователя (в admin или текущей БД)
db = db.getSiblingDB("admin");
print(">>> INIT SCRIPT STARTED");

db.createUser({
	user: "articleService",
	pwd: "{{ .Env.ARTICLE_SERVICE_PASSWORD }}",
	roles: [
		{
			role: "readWrite",
			db: "article-storage"
		}
	]
});

// создаем БД и коллекцию
db = db.getSiblingDB("article-storage");

// создание коллекции
db.createCollection("articles");

// можно сразу добавить документ (создаст БД автоматически)
// db.users.insertOne({
// 	name: "Admin",
// 	createdAt: new Date()
// });