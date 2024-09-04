
# HETIC Projet Final: “ Tech-Race “

<img src="https://img.shields.io/badge/golang-%5E1.22-blue">

<div align="center">
   <img src="https://github.com/NasssDev/Tech-Race/assets/167258734/8022059e-d34b-422f-9010-bf8d8fdd7132" alt="image" width="300" height="200"/>
</div>

## Description du Projet

Le projet final HETIC, "Tech Race", vise à développer une application mobile permettant de contrôler un véhicule à distance via un réseau sans fil. L'application offre aux utilisateurs la possibilité de piloter le véhicule, d'accéder aux données de télémétrie en temps réel et de participer à des courses autonomes.

Ce projet se compose de plusieurs repositories toutes hébergés sur Github : 
- [API Tech Race](https://github.com/NasssDev/Tech-Race)
- [App Mobile](https://github.com/Hetic-Team/tech_race_8_2024)
- [Site web de partage des vidéos](https://site-a-venir)
- [Programme de la voiture](https://github.com/ExploryKod/freenove_esp32_wrover)
- [Modèle de la voiture](https://www.amazon.fr/Freenove-ESP32-WROVER-Contained-Compatible-Expressions/dp/B08X6PTQFM/ref=sr_1_5?__mk_fr_FR=%C3%85M%C3%85%C5%BD%C3%95%C3%91&crid=1NFTVTE5M400B&dib=eyJ2IjoiMSJ9.ouyBflLDqHVkfViARMLD6Bn9gOI47kLGrM-5LMAbtJPAUgPogSQ1tQyH60VxNGSHTf-JIYDTkVL4RJ2a7-L92dQ5aqD8IliDd4MzLvffNmw65QxSItZh_qi-vPHXgzjBhvcW8Vy00EckrayFx_47OCj3W4K6Y1W0jHZgIDF7DAvRTI9XcC7oRK8T9xeUORe35q6RJ29TNUuhLCcN5fXl-WqLhsgNb2JA0XzHwnqwHaBBwj-xZ77ohEfVpUYfdyOMWf1wO01Fa42MzKl0b-UGD6PwYD-kBCJYQS3J9twWSGs.OrlAkZRIvlaYtQ2-9pywcADOLR7VY4iRx_9Ps1DkMnk&dib_tag=se&keywords=esp32+car&qid=1715602634&sprefix=esp+32+car,aps,125&sr=8-5)


## Equipe Backend (ce repository) : 

- [Amaury FRANSSEN](https://github.com/ExploryKod) 
- [Nassim AISSAOUI](https://github.com/NasssDev)
- [Justin LELUC](https://github.com/Jykiin)

## Equipe frontend : 
- [Reewaz Maskey](https://github.com/reewaz001)
- [Alexandre VISAGE](https://github.com/Aleex470)
- [Khalifa boubacar DIONE](https://github.com/khalifadione)
- [Achraf CHARDOUDI](https://github.com/Achkey)


## Objectif

L'objectif principal est de créer une interface intuitive pour les pilotes afin de contrôler le véhicule et de visualiser ses performances. Deux niveaux de courses sont organisés : l'un où les pilotes contrôlent directement le véhicule depuis l'application mobile, et l'autre où le véhicule doit naviguer de manière autonome en suivant une ligne au sol.

## Fonctionnalités Principales

- Contrôle à distance du véhicule via l'application mobile.
- Visualisation en temps réel des données de télémétrie (vitesse, orientation, etc.).
- Mode suiveur de ligne autonome pour les courses sans intervention humaine.

## ⚙️ Configuration locale de l'API Backend (ce repository)

#### Base générale
1. Cloner ce repository 
2. Avoir installé Golang sur sa machine locale [Installation de Go](https://go.dev/doc/install)
3. Avoir docker sur sa machine locale et avoir installé l'image docker de mosquitto : [Cliquez ici](https://github.com/ExploryKod/mosquitto-docker) 
4. Avoir postgresSQL (optionnel car on peux passer via docker)
5. Installer MakeFile si vous êtes sur une platforme qui ne l'a pas nativement 

#### Pour faire fonctionner le service de vidéos
4. S'inscrire gratuitement sur Cloudinary : [S'inscrire sur Cloudinary](https://cloudinary.com/)
5. Récupérer son cloudinary ID et le cloudinary URL depuis son compte
6. Configurer le service cloudinarace (service cloudinary de Tech Race) : 
**Depuis la racine, se rendre dans `pkg/other/cloudinary/.env` :** 
```
CLOUDINARY_ID=mon-id-cloudinary-présente-sur-mon-compte
CLOUDINARY_URL=mon-url-cloudinary-présente-sur-mon-compte
GOOS=linux (ou d'autres os si vous n'avez pas linux) > doc de Golang
GOARCH=amd64 (vérifier aussi que c'est bien la bonne architecture pour vous aussi) > doc de Golang
```

#### 
4. Configurer le `.env` du projet à la racine : 
``` 
PORT_VIDEO=7000
IP=127.0.0.1
BOUNDARY=--123456789000000000000987654321
```

5. Aller à la racine du projet et lancer ces commandes : 
Vérifier que les ports 1883 (mqtt), 8083 (cloudinarace), 9000, 8888 (pgAdmin), 8889 (adminer) et 5432 (postgres) sont libres.
- Lancer docker : `docker compose up -d`
- Installer les dépendances : `go mod download` et `go mod tidy`
- Lancer le projet : `go run cmd/api/main.go`
- Lancer cloudinarace pour la gestion vidéo (si besoin): `make cloudinarace`
Si vous n'avez pas Make, lancez cloudinarace via : `cd pkg/other/cloudinary && go run main/main.go --port=8083`

6. Configurer la base de donnée: 
- Se rendre sur [adminer](http://localhost:8089) ou [pgAdmin](http://localhost:8888)
- Regarder les credentials présents dans le fichier `docker-compose.yaml` 
- Choisir PostgresSQL et utiliser ces credentials : 

```
Bdd : tech_race  
User : root
Sur PgAdmin, default email : tech@race.com
Mot de passe : password
serveur : db 
```

- Import le dump de la base de donnée

## Matériel et Logiciel Utilisés

### Matériel
- Architecture ESP32
- 4 roues motrices pilotées indépendamment
- Connexion sans fil wifi/bluetooth
- Capteur de distance
- Capteur de suivi de ligne
- Caméra embarquée

### Logiciel
- APIs de contrôle pour les roues et les capteurs
- Données de télémétrie fournies par les capteurs
- Communication en temps réel entre l'application mobile et le véhicule

## Instructions d'Installation

1. Clonez le dépôt GitHub vers votre machine locale.
2. Assurez-vous d'avoir les outils de développement nécessaires installés (par exemple, Node.js, npm).
3. Installez les dépendances du projet en exécutant `npm install`.
4. Démarrez l'application en exécutant `npm start`.

## Contribuer

Les contributions sont les bienvenues ! Pour contribuer à ce projet, veuillez suivre ces étapes :
1. Forker le projet
2. Créer une branche pour votre fonctionnalité (`git checkout -b feature/AmazingFeature`)
3. Commiter vos modifications (`git commit -m 'Ajouter une fonctionnalité incroyable'`)
4. Pousser vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

## Licence

Ce projet est sous licence MIT. Pour plus d'informations, veuillez consulter le fichier [LICENSE](LICENSE).
