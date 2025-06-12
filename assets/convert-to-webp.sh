#!/bin/bash

echo "🔁 Mulai ulang konversi semua gambar ke WebP (paksa overwrite)..."
echo "📂 Direktori kerja: $(pwd)"
echo

find . -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) | while read img; do
  webp_name="${img%.*}.webp"

  echo "🖼️  Mengonversi ulang: $img"
  cwebp -q 60 -m 6 -pass 10 -resize 1280 0 "$img" -o "$webp_name" > /dev/null 2>&1

  if [ $? -eq 0 ]; then
    original_size=$(du -h "$img" | cut -f1)
    webp_size=$(du -h "$webp_name" | cut -f1)
    echo "✅ Overwrite: $webp_name (dari $original_size ke $webp_size)"

    rm "$img"
    echo "🗑️  File asli dihapus: $img"
  else
    echo "❌ Gagal mengonversi $img"
  fi

  echo
done

echo "🎉 Semua gambar berhasil diperbarui ke WebP!"
