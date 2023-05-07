#!/usr/bin/env python3
import json
import os
import sys
from PIL import Image, ImageFilter, ImageDraw, ImageFont


OVERLAP = 0.16

assets_folder = "images"
if os.environ["ASSETS_PATH"] != "":
    assets_folder = os.path.abspath(os.environ["ASSETS_PATH"])
print(f"Loading images from {assets_folder=}")
genshin_font = ImageFont.truetype(
    os.path.join(assets_folder, "fonts/genshin_font.ttf"), 120)

def get_data() -> dict:
    return json.load(sys.stdin)

def open_image(fp):
    try:
        return Image.open(fp)
    except Exception:
        try:
            return Image.open(os.path.join(os.path.dirname(fp), "default.png"))
        except Exception as e:
            print(e)
            return Image.new("RGBA", (256, 256))
data = get_data()
incomplete_chars = []
if "incomplete_chars" in data.keys():
    incomplete_chars.extend(data["incomplete_chars"])
chars = data["char_details"]
# print(chars[0])
names = [chars[x]["name"] for x in range(len(chars))]
weapons = [chars[x]["weapon"] for x in range(len(chars))]
artifacts = [chars[x]["sets"] for x in range(len(chars))]

char_image_shapes = []
imgs = []
new_image_width = 900
new_image_height = 422
for name in names:
    char_img = open_image(os.path.join(assets_folder, f"avatar/{name}.png"))
    char_img = char_img.resize((256,256),Image.Resampling.BICUBIC)
    if name in incomplete_chars:
        text_img = Image.new("RGBA", (256,256), (0, 0, 0, 0))
        text = ImageDraw.Draw(text_img)
        text.text((256/2,256/2), "WIP", font=genshin_font, fill=(0,0,0), anchor="mm")
        text_img = text_img.rotate(60, resample=Image.Resampling.BILINEAR, expand=False, center=(256/2+20,256/2))
        shadow = Image.new("RGBA", text_img.size, (255, 255, 255, 255))
        alpha = text_img.split()[-1]
        shadow.putalpha(alpha)
        shadow = shadow.filter(ImageFilter.MaxFilter(7))
        shadow = shadow.filter(ImageFilter.GaussianBlur(1))
        shadow.alpha_composite(text_img)
        char_img.alpha_composite(shadow)

    imgs.append(char_img)
    char_image_shapes.append((256, 256))

base_img = Image.new("RGBA", (new_image_width, new_image_height))
location = [[0, 0] for _ in range(len(imgs))]
for i in range(len(imgs)-1):
    location[i+1][0] = location[i][0] + \
        int(char_image_shapes[i][0] * (1-OVERLAP))
# print(location)

# for i in range(len(imgs)):
#     shadow = Image.new("RGBA", char_image_shapes[i], (255, 255, 255, 255))
#     alpha = imgs[i].split()[-1]
#     shadow.putalpha(alpha)
#     shadow = shadow.filter(ImageFilter.MaxFilter(5))
#     shadow.alpha_composite(imgs[i])
#     imgs[i] = shadow

# characters
for i in range(len(imgs)-1, -1, -1):
    img = imgs[i]
    base_img.alpha_composite(img, tuple(location[i]))

# weapons
# print(weapons)
weapon_size = (180, 180)
imgs: list[Image.Image] = []
weapon_image_shapes = []
for weapon in weapons:

    imgs.append(open_image(os.path.join(
        assets_folder, f"weapons/{weapon['name']}.png")))
    width, height = imgs[-1].size
    imgs[-1] = imgs[-1].resize(weapon_size)
    imgs[-1] = imgs[-1].crop(imgs[-1].getbbox())
    weapon_image_shapes.append(imgs[-1].size)

weapon_img = Image.new("RGBA", base_img.size, (0, 0, 0, 0))
for i in range(len(imgs)):
    img = imgs[i]
    weapon_img.alpha_composite(img,  (location[i][0] + char_image_shapes[i][0] - int(
        weapon_image_shapes[i][0] * 0.95)-10, char_image_shapes[i][1] - weapon_image_shapes[i][1] + 20))

shadow = Image.new("RGBA", weapon_img.size, (255, 255, 255, 255))
alpha = weapon_img.split()[-1]
shadow.putalpha(alpha)
shadow = shadow.filter(ImageFilter.MaxFilter(5))
shadow.alpha_composite(weapon_img)
weapon_img = shadow

shadow = Image.new("RGBA", weapon_img.size, (0, 0, 0, 255))
alpha = weapon_img.split()[-1]
shadow.putalpha(alpha)
shadow = shadow.filter(ImageFilter.MaxFilter(7))
shadow = shadow.filter(ImageFilter.GaussianBlur(2))
shadow.alpha_composite(weapon_img)
weapon_img = shadow

base_img.alpha_composite(weapon_img,  (0, 0))

imgs: list[Image.Image] = []
artifact_image_shapes = []
ARITFACT_SIZE = (100, 100)
for arti in artifacts:
    arti = {key: val for key, val in arti.items() if val >= 2}

    sets = list(arti.keys())
    total_sets = len(sets)
    # print(sets)
    if total_sets == 1:
        if arti[sets[0]] >= 4:
            imgs.append(open_image(os.path.join(
                assets_folder, f"artifacts/{sets[0]}_flower.png")))
            imgs[-1] = imgs[-1].resize(ARITFACT_SIZE)
        else:
            img0 = open_image(os.path.join(
                assets_folder, f"artifacts/{sets[0]}_flower.png"))
            img0 = img0.resize(ARITFACT_SIZE)
            img0 = img0.crop((0, 0, ARITFACT_SIZE[0]//2, ARITFACT_SIZE[1]))
            dst = Image.new("RGBA", ARITFACT_SIZE, (0, 0, 0, 0))
            dst.paste(img0, (0, 0))
            dst_draw = ImageDraw.Draw(dst)
            dst_draw.line(
                (ARITFACT_SIZE[0]//2, 0, ARITFACT_SIZE[0]//2, ARITFACT_SIZE[1]), fill=0, width=4)
            imgs.append(dst)
    elif total_sets == 2:
        img0 = open_image(os.path.join(
            assets_folder, f"artifacts/{sets[0]}_flower.png"))
        img0 = img0.resize(ARITFACT_SIZE)
        img0 = img0.crop((0, 0, ARITFACT_SIZE[0]//2, ARITFACT_SIZE[1]))

        img1 = open_image(os.path.join(
            assets_folder, f"artifacts/{sets[1]}_flower.png"))
        img1 = img1.resize(ARITFACT_SIZE)
        img1 = img1.crop(
            (ARITFACT_SIZE[0]//2, 0, ARITFACT_SIZE[0], ARITFACT_SIZE[1]))

        dst = Image.new("RGBA", ARITFACT_SIZE, (0, 0, 0, 0))
        dst.paste(img0, (0, 0))
        dst.paste(img1, (img0.width, 0))
        dst_draw = ImageDraw.Draw(dst)
        dst_draw.line(
            (ARITFACT_SIZE[0]//2, 0, ARITFACT_SIZE[0]//2, ARITFACT_SIZE[1]), fill=0, width=4)
        imgs.append(dst)
    else:
        imgs.append(Image.new("RGBA", ARITFACT_SIZE, (0, 0, 0, 0)))
    artifact_image_shapes.append(ARITFACT_SIZE)

for i in range(len(imgs)):
    shadow = Image.new("RGBA", artifact_image_shapes[i], (255, 255, 255, 255))
    alpha = imgs[i].split()[-1]
    shadow.putalpha(alpha)
    shadow = shadow.filter(ImageFilter.MaxFilter(5))
    shadow.alpha_composite(imgs[i])
    imgs[i] = shadow

for i in range(len(imgs)):
    shadow = Image.new("RGBA", artifact_image_shapes[i], (0, 0, 0, 255))
    alpha = imgs[i].split()[-1]
    shadow.putalpha(alpha)
    shadow = shadow.filter(ImageFilter.MaxFilter(7))
    shadow = shadow.filter(ImageFilter.GaussianBlur(2))
    shadow.alpha_composite(imgs[i])
    imgs[i] = shadow

for i in range(len(imgs)):
    img = imgs[i]
    base_img.alpha_composite(img,  (location[i][0] + char_image_shapes[i][0] -
                                    artifact_image_shapes[i][0]-10, char_image_shapes[i][1] - artifact_image_shapes[i][1] + 40))
    # base_img.alpha_composite(img,  (location[i][0] + 50, new_image_height - artifact_image_shapes[i][1] + 20))
    # base_img.alpha_composite(img,  (location[i][0] + 20, 0))

blue = (102, 170, 206, 255)
purple = (154, 112, 197, 255)
gold = (217, 170, 91, 255)
white = (255, 255, 255, 255)
genshin_font = ImageFont.truetype(
    os.path.join(assets_folder, "fonts/genshin_font.ttf"), 30)
text_img = Image.new("RGBA", base_img.size, (0, 0, 0, 0))
text = ImageDraw.Draw(text_img)

for i in range(len(chars)):
    char = chars[i]
    cons = char["cons"]
    ref = char["weapon"]["refine"]

    # text.text((location[i][0] + 30, new_image_height - 30), f"C{cons}", font = genshin_font, fill = blue)
    # xDescPxl = text.textsize(f"C{cons}", font= genshin_font)[0]
    # text.text((location[i][0] + 30 + xDescPxl, new_image_height - 30), f"R{ref}", font = genshin_font, fill = gold)
    text.text((location[i][0] + 50, 5),
              f"C{cons}", font=genshin_font, fill=gold)
    xDescPxl = text.textsize(f"C{cons}", font=genshin_font)[0]
    text.text((location[i][0] + 50 + xDescPxl, 5),
              f"R{ref}", font=genshin_font, fill=purple)

# duration
# genshin_font_28px = ImageFont.truetype("genshin_font.ttf", 28)
dps = data["dps"]
info = f"""
Total DPS: {dps['mean']:5.0f} to {data['num_targets']} target{'s' if data['num_targets'] > 1 else ''} (Avg. Per Target {dps['mean']/data['num_targets']:5.0f})
DPS min / max / stddev: {dps['min']:.0f} / {dps['max']:.0f} / {dps['sd']:.0f}
{data['sim_duration']['mean']:.2f}s combat time. {data['iter']} iteration in {(data['runtime']/1e9):.3f}s
"""
text.text((6, char_image_shapes[0][1]), info,
          font=genshin_font, fill=white, spacing=10)

# shadow = Image.new("RGBA", text_img.size, (255, 255, 255, 255))
# alpha = text_img.split()[-1]
# shadow.putalpha(alpha)
# shadow = shadow.filter(ImageFilter.MaxFilter(3))
# shadow.alpha_composite(text_img)
# text_img = shadow

shadow = Image.new("RGBA", text_img.size, (0, 0, 0, 255))
alpha = text_img.split()[-1]
shadow.putalpha(alpha)
shadow = shadow.filter(ImageFilter.MaxFilter(7))
shadow = shadow.filter(ImageFilter.GaussianBlur(1))
shadow.alpha_composite(text_img)
text_img = shadow

base_img.alpha_composite(text_img)

base_img = base_img.resize(map(lambda x: int(x*0.6), base_img.size))

output_filename = "output.png"
if len(sys.argv) > 1:
    output_filename = sys.argv[1]
    if not output_filename.endswith(".png"):
        output_filename += (".png")
base_img.save(output_filename)
# base_img.show()
