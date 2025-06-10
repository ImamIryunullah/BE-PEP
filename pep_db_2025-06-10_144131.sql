-- MySQL dump 10.13  Distrib 8.0.33, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: pep_db
-- ------------------------------------------------------
-- Server version	8.4.3

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `berita`
--

DROP TABLE IF EXISTS `berita`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `berita` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `judul` varchar(255) NOT NULL,
  `subtitle` varchar(255) DEFAULT NULL,
  `tanggal` datetime(3) DEFAULT NULL,
  `penulis` varchar(100) NOT NULL,
  `isi` text NOT NULL,
  `foto` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_berita_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `berita`
--

/*!40000 ALTER TABLE `berita` DISABLE KEYS */;
INSERT INTO `berita` VALUES (1,'2025-06-09 15:07:28.993','2025-06-09 15:31:04.020',NULL,'edittttttttttttttttttttt','fwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww','2025-06-09 07:00:00.000','wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww','feeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee','1749456448_esport.png'),(2,'2025-06-09 15:07:49.105','2025-06-09 15:07:49.105',NULL,'fwafwawfawfwafwafwafwad','fwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww','2025-06-09 07:00:00.000','wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww','fwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwffffwwwwww','1749456469_tesnimeja.png'),(3,'2025-06-09 15:10:54.591','2025-06-09 15:10:54.591',NULL,'fwafwawfawfwafwafwafwad','fwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww','2025-06-09 07:00:00.000','wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww','wfafafffffffffffffffffffffffffffffffffffffffffffffffffffffffff','1749456654_voli.png'),(4,'2025-06-09 15:11:46.389','2025-06-09 15:11:46.389','2025-06-09 15:31:10.804','fgeegegegeg','gggggggggggggggggggggggggggggggggg','2025-06-09 07:00:00.000','ggggggg','gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg','1749456706_lari.png'),(5,'2025-06-09 17:23:52.853','2025-06-09 17:23:52.853',NULL,'berita baruuuuuuu','newwww','2025-06-09 07:00:00.000','TEST','hudddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd','1749464632_Hero.png');
/*!40000 ALTER TABLE `berita` ENABLE KEYS */;

--
-- Table structure for table `daftar_users`
--

DROP TABLE IF EXISTS `daftar_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `daftar_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `aset` varchar(100) NOT NULL,
  `provinsi` varchar(100) NOT NULL,
  `foto` varchar(255) DEFAULT NULL,
  `foto_path` varchar(500) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `last_login` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_daftar_users_email` (`email`),
  KEY `idx_daftar_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `daftar_users`
--

/*!40000 ALTER TABLE `daftar_users` DISABLE KEYS */;
INSERT INTO `daftar_users` VALUES (1,'2025-06-08 02:04:30.123','2025-06-08 02:04:30.123',NULL,'test@gmail.com','$2a$10$qnfpuwtUd5NFQXcFWedr/Oz1n6C1uAfbNGRkhco4OpKjr8IdYEt8i','Aset1','JawaTimur','1749323070_esport.png','uploads\\1749323070_esport.png',1,NULL);
/*!40000 ALTER TABLE `daftar_users` ENABLE KEYS */;

--
-- Table structure for table `funruns`
--

DROP TABLE IF EXISTS `funruns`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `funruns` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `nama` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `kontingen` varchar(100) NOT NULL,
  `size` varchar(10) NOT NULL,
  `status` varchar(191) DEFAULT 'pending',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_funruns_email` (`email`),
  KEY `idx_funruns_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `funruns`
--

/*!40000 ALTER TABLE `funruns` DISABLE KEYS */;
INSERT INTO `funruns` VALUES (1,'2025-06-09 02:55:08.693','2025-06-09 03:45:34.168',NULL,'Imam','m.imam.2205356@students.um.ac.id','Aset 1','XS','approved'),(2,'2025-06-09 02:56:00.508','2025-06-09 03:45:34.189',NULL,'Pekerja','test@gmail.com','Aset 5','XXL','approved'),(3,'2025-06-09 02:56:00.516','2025-06-09 03:45:34.191',NULL,'wfwfawf','fwhfwf@gmail.com','Aset 4','S','approved'),(4,'2025-06-09 15:06:35.652','2025-06-09 15:06:55.236',NULL,'ddsssd','tatyahanum10@gmail.com','Aset 2','S','approved'),(5,'2025-06-09 17:25:08.322','2025-06-09 17:25:08.322',NULL,'murulaki','murudi@gmail.com','Aset 2','S','pending');
/*!40000 ALTER TABLE `funruns` ENABLE KEYS */;

--
-- Table structure for table `knockout_matches`
--

DROP TABLE IF EXISTS `knockout_matches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `knockout_matches` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `kategori` varchar(100) NOT NULL,
  `sub_kategori` varchar(100) NOT NULL,
  `tim1` varchar(255) NOT NULL,
  `tim2` varchar(255) NOT NULL,
  `hasil` varchar(100) NOT NULL,
  `tahap` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_knockout_matches_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `knockout_matches`
--

/*!40000 ALTER TABLE `knockout_matches` DISABLE KEYS */;
/*!40000 ALTER TABLE `knockout_matches` ENABLE KEYS */;

--
-- Table structure for table `participant_registrations`
--

DROP TABLE IF EXISTS `participant_registrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `participant_registrations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `nama_lengkap` longtext,
  `email` longtext,
  `no_telepon` longtext,
  `jenis_peserta` longtext,
  `cabang_olahraga` longtext,
  `wilayah_kerja` longtext,
  `media_sosial` longtext,
  `catatan` longtext,
  `waktu_daftar` datetime(3) DEFAULT NULL,
  `ktp` longtext,
  `id_card` longtext,
  `bpjs` longtext,
  `pas_foto` longtext,
  `surat_pernyataan` longtext,
  `surat_layak_bertanding` longtext,
  `form_prq` longtext,
  `surat_keterangan_kerja` longtext,
  `kontrak_kerja` longtext,
  `sertifikat_bst` longtext,
  `status` varchar(191) DEFAULT 'pending',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_participant_registrations_user_id` (`user_id`),
  KEY `idx_participant_registrations_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_daftar_users_registrations` FOREIGN KEY (`user_id`) REFERENCES `daftar_users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `participant_registrations`
--

/*!40000 ALTER TABLE `participant_registrations` DISABLE KEYS */;
INSERT INTO `participant_registrations` VALUES (1,1,'M. IMAM IRYUNULLAH','','0859171293969','Pekerja','Tenis Meja','Jakarta','','dqdqdqdqd','2025-06-08 06:04:57.672','20250608060457_ktp.png','20250608060457_id_card.png','20250608060457_bpjs.png','20250608060457_pas_foto.png','20250608060457_surat_pernyataan.png','20250608060457_surat_layak_bertanding.png','20250608060457_form_prq.png','','','','approved','2025-06-08 06:04:57.678','2025-06-08 21:56:24.140',NULL),(2,1,'dwdwfwafawf','','0859171293969','Mitra','Voli','Jakarta','https://otakudesu.cloud/anime/hky-season-2-sub-indo/','dqdqd','2025-06-08 21:57:35.206','20250608215735_ktp.png','20250608215735_id_card.png','20250608215735_bpjs.png','20250608215735_pas_foto.png','20250608215735_surat_pernyataan.png','20250608215735_surat_layak_bertanding.png','20250608215735_form_prq.png','20250608215735_surat_keterangan_kerja.png','20250608215735_kontrak_kerja.png','20250608215735_sertifikat_bst.png','rejected','2025-06-08 21:57:35.212','2025-06-08 21:57:43.919',NULL),(3,1,'Muru Rabangodufddw','','085736352624','Pekerja','Basket','Malang','','OOOOO','2025-06-09 08:34:14.897','20250609083414_ktp.png','20250609083414_id_card.png','20250609083414_bpjs.png','20250609083414_pas_foto.png','20250609083414_surat_pernyataan.png','20250609083414_surat_layak_bertanding.png','20250609083414_form_prq.png','','','','rejected','2025-06-09 08:34:14.906','2025-06-09 14:53:07.461',NULL);
/*!40000 ALTER TABLE `participant_registrations` ENABLE KEYS */;

--
-- Dumping routines for database 'pep_db'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-06-10 14:42:03
