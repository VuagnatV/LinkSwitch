LinkSwitch 

	Un backend pour pour réduire les urls

Requirements

	Une base de donnée MongoDB
	(Optionnel) Environement GoLang

Lancement 

	Pour lancer le server :
	- go run main.go

	Pour lancer les tests:
	- go test -v -cover ./...
 
Environnement 

	Port : ":8000"
  
	MongoDB : "mongodb://localhost:27017"

Routes

	POST : / : {long : string } : Prend un objet avec url et le renvoie avec une version racourcie
  
	GET : /{shortURL} :  : Prend une url raccourcie et redirige vers l'url longue, error si l'url a expiré
  
	GET : /stats/ :  : Retourne le nombre d'urls raccourcies
  
	GET : /stats/{shortURL} :  : Prend une url raccourcie et retourne le nombre de fois que l'url à été clické

 	PUT : /{shortURL} : {short : string} : Prend une url raccourcie et un objet avec une url courte et modifie l'url courte 
