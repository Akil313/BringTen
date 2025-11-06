import os
import xml.etree.ElementTree as ET

def remove_svg_dimensions(file_path):
    ET.register_namespace('', "http://www.w3.org/2000/svg")  # Prevent xmlns duplication
    try:
        tree = ET.parse(file_path)
        root = tree.getroot()

        # Remove width and height only from the <svg> tag
        for attr in ["width", "height"]:
            if attr in root.attrib:
                del root.attrib[attr]

        tree.write(file_path, encoding='utf-8', xml_declaration=True)
        print(f"✅ Cleaned: {file_path}")
    except Exception as e:
        print(f"❌ Error with {file_path}: {e}")

def clean_all_svgs(folder_path):
    for root_dir, _, files in os.walk(folder_path):
        for file in files:
            if file.endswith(".svg"):
                file_path = os.path.join(root_dir, file)
                remove_svg_dimensions(file_path)

if __name__ == "__main__":
    svg_folder = "../bring-ten/src/lib/images/cards"  # 🔁 Change to your folder path
    clean_all_svgs(svg_folder)
