
# HETC Web: Tech Race

## Description du Projet

Le projet HETC Web: Tech Race vise à développer une application mobile permettant de contrôler un véhicule à distance via un réseau sans fil. L'application offre aux utilisateurs la possibilité de piloter le véhicule, d'accéder aux données de télémétrie en temps réel et de participer à des courses autonomes.

## Objectif

L'objectif principal est de créer une interface intuitive pour les pilotes afin de contrôler le véhicule et de visualiser ses performances. Deux niveaux de courses sont organisés : l'un où les pilotes contrôlent directement le véhicule depuis l'application mobile, et l'autre où le véhicule doit naviguer de manière autonome en suivant une ligne au sol.

## Fonctionnalités Principales

- Contrôle à distance du véhicule via l'application mobile.
- Visualisation en temps réel des données de télémétrie (vitesse, orientation, etc.).
- Mode suiveur de ligne autonome pour les courses sans intervention humaine.

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
