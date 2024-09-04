
# HETIC Projet Final: “ Tech-Race “

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

- [Amaury](https://github.com/ExploryKod) 
- [Nassim](https://github.com/NasssDev)
- [Justin](https://github.com/Jykiin)

## Objectif

L'objectif principal est de créer une interface intuitive pour les pilotes afin de contrôler le véhicule et de visualiser ses performances. Deux niveaux de courses sont organisés : l'un où les pilotes contrôlent directement le véhicule depuis l'application mobile, et l'autre où le véhicule doit naviguer de manière autonome en suivant une ligne au sol.

## Fonctionnalités Principales

- Contrôle à distance du véhicule via l'application mobile.
- Visualisation en temps réel des données de télémétrie (vitesse, orientation, etc.).
- Mode suiveur de ligne autonome pour les courses sans intervention humaine.

## ⚙️ Configuration locale du projet

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

## Membres du projet

1. Justin LELUC (@)
2. Reewaz MASKEY (@)
3. Achraf CHARDOUDI (@)
4. Amaury FRANSSEN (@)
5. Alexandre VISAGE (@)
6. Nassim AISSAOUI (@)
7. Khalifa boubacar DIONE (@)

## Contribuer

Les contributions sont les bienvenues ! Pour contribuer à ce projet, veuillez suivre ces étapes :
1. Forker le projet
2. Créer une branche pour votre fonctionnalité (`git checkout -b feature/AmazingFeature`)
3. Commiter vos modifications (`git commit -m 'Ajouter une fonctionnalité incroyable'`)
4. Pousser vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

## Licence

Ce projet est sous licence MIT. Pour plus d'informations, veuillez consulter le fichier [LICENSE](LICENSE).
