import os
import sys
import win32com.client

def create_shortcut(exe_path, shortcut_name=None):
    if not os.path.exists(exe_path):
        print(f"Erreur : Le fichier {exe_path} n'existe pas.")
        return

    # Définir le nom du raccourci
    if shortcut_name is None:
        shortcut_name = os.path.splitext(os.path.basename(exe_path))[0]

    # Déterminer le bureau de l'utilisateur
    desktop = os.path.join(os.path.expanduser("~"), "Desktop")
    shortcut_path = os.path.join(desktop, f"{shortcut_name}.lnk")

    # Récupérer le répertoire contenant l'exécutable
    work_dir = os.path.dirname(exe_path)

    # Créer un raccourci avec le bon chemin et répertoire de travail
    iconPath = "C:/Users/rayou/Documents/GitHub/Sudoku-GO/details/assets/icone.ico"
    shell = win32com.client.Dispatch("WScript.Shell")
    shortcut = shell.CreateShortcut(shortcut_path)
    shortcut.TargetPath = exe_path
    shortcut.WorkingDirectory = work_dir  # Définit le dossier de travail
    shortcut.IconLocation = iconPath  # Utilise l'icône
    shortcut.Save()

    print(f"Raccourci '{shortcut_name}.lnk' créé sur le bureau.")

# Exemple d'utilisation :
exe_file = r"C:/Users/rayou/Documents/GitHub/Sudoku-GO/sudoku.exe"
create_shortcut(exe_file)
