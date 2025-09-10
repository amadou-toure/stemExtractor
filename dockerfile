# Utilise une image Python officielle légère
FROM python:3.10-bullseye

# Met à jour et installe ffmpeg (pour lire/écrire les fichiers audio)
RUN apt update && apt install -y ffmpeg git libsndfile1 && rm -rf /var/lib/apt/lists/*

# Crée un dossier de travail
WORKDIR /app

# Copie les fichiers requirements
COPY requirements.txt .

# Installe les dépendances Python
RUN pip install --upgrade pip && pip install -r requirements.txt

# Copie le reste de l’application (ex: app.py)
COPY . .

# Déclare le point d’entrée de l’application
CMD ["demucs", "--help"]