
# Cloudinarace 

## Gestion de vidéos sur cloudinary

### Cloudinarace et Tech Race

Le projet HETIC "Tech Race" vise à développer une application mobile permettant de contrôler un véhicule à distance via un réseau sans fil. 
L'application offre aux utilisateurs la possibilité de piloter le véhicule, d'accéder aux données de télémétrie en temps réel et de participer à des courses autonomes.

"Tech Race" ici désigne aussi la partie backend écrite en Goland auquel ce package se rattache.

Dans le cadre de ce projet, nous créons des packages indépendants qui ont vocation première de fonctionner avec le coeur de l'application "Tech Race". 

Ce package est l'un d'entre-eux.

### Objectif

L'objectif principal est de traiter et gérer l'export de vidéos prise par la caméra de la voiture sur Cloudinary.

### Instructions d'Installation

1. Clonez le dépôt GitHub vers votre machine locale.
2. Assurez-vous d'avoir les outils de développement nécessaires installés (Go sur votre machine locale).
3. Installez des conteneurs docker : `make cloud-docker`.
4. Démarrez l'application en exécutant : <br>`make start-cloudinarace PORT=votre-port`.
5. Se rendre sur la page d'accueil pour suivre les explications en vue d'utiliser ce package

### Licence

Ce projet est sous licence MIT. Pour plus d'informations, veuillez consulter le fichier [LICENSE](LICENSE).
