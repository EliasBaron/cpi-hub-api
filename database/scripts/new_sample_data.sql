--CPI HUB SAMPLE DATA
-- This script populates the database with sample data for users, spaces, posts, and comments.
-- It includes 15 users, 20 spaces, 70 posts, and 140 comments to simulate a realistic environment for testing and development.
-- =========================

-- Insert sample users
INSERT INTO users (name, last_name, email, password, image, created_at, updated_at) VALUES
                                                                                        ('Valentin', 'Ferreyra', 'valentin.ferreyra@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/08/fc/90/08fc9089f2cbd576edc647e075f5eb0a.jpg', NOW(), NOW()),
                                                                                        ('Elias', 'Baron', 'elias.baron@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/1200x/78/ec/04/78ec04232efb3c34f4a6fa57a1c62f0a.jpg', NOW(), NOW()),
                                                                                        ('Juanma', 'Sanchez Diaz', 'juanma.sanchez.diaz@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/0b/1e/1d/0b1e1dfa6a93149a2277f400a97e3ce8.jpg', NOW(), NOW()),
                                                                                        ('Diego', 'Kippes', 'diego.kippes@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/33/72/f5/3372f5e7d7642babfffb889e5fd90133.jpg', NOW(), NOW()),
                                                                                        ('Aaron', 'Gutierrez', 'aaron.gutierrez@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/1200x/8f/91/23/8f91233d4ed1a0497214cd17c04d624f.jpg', NOW(), NOW()),
                                                                                        ('Santiago', 'Abregu', 'santiago.abregu@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/8d/04/96/8d04969b987b7f22738af69d5854103b.jpg', NOW(), NOW()),
                                                                                        ('Fabian', 'Frangella', 'fabian.frangella@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/4b/42/fb/4b42fbe594dae2cacfa200c3d2b54f15.jpg', NOW(), NOW()),
                                                                                        ('Matias', 'Aduco', 'matias.aduco@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/06/5a/40/065a405b2259e863cfcb7570663ff7cb.jpg', NOW(), NOW()),
                                                                                        ('Gonzalo', 'Bender', 'gonzalo.bender@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/d8/9a/b6/d89ab6d59778d9602ae0a164802c5322.jpg', NOW(), NOW()),
                                                                                        ('Tobias', 'Torres', 'tobias.torres@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/50/c0/08/50c008dc56e5f9fe49a2912d35fc964b.jpg', NOW(), NOW()),
                                                                                        ('Carlos', 'Rivero', 'carlos.rivero@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/dd/e5/e4/dde5e4e9bb3e8b5835a9bae5a2be93c5.jpg', NOW(), NOW()),
                                                                                        ('Susan', 'Rosito', 'susan.rosito@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/736x/30/96/36/309636caa74f83f80fb20c52b401db58.jpg', NOW(), NOW()),
                                                                                        ('Fernando', 'Dodino', 'fernando.dodino@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://i.pinimg.com/1200x/67/1a/1a/671a1a21eb014736397250a44d1469db.jpg', NOW(), NOW()),
                                                                                        ('Feche', 'Romero', 'feche.romero@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'https://ar.pinterest.com/pin/26106872836453625/', NOW(), NOW());

-- Insert sample spaces
INSERT INTO spaces (name, description, created_by, updated_by, created_at, updated_at) VALUES
                                                                                           ('Club de Fotografía 📸', 'Comparte tus mejores fotos y técnicas fotográficas', 1, 1, NOW(), NOW()),
                                                                                           ('Grupo Fitness 💪', 'Rutinas de ejercicio, consejos de nutrición y motivación deportiva', 2, 2, NOW(), NOW()),
                                                                                           ('Cine y Series 🎬', 'Reseñas de películas, series y recomendaciones cinematográficas', 3, 3, NOW(), NOW()),
                                                                                           ('Emprendimientos 🚀', 'Ideas de negocio, startups y consejos para emprendedores', 4, 4, NOW(), NOW()),
                                                                                           ('Grupo de Idiomas 🌍', 'Aprende idiomas, practica conversación y comparte recursos', 5, 5, NOW(), NOW()),
                                                                                           ('Grupo de Medio Ambiente 🌱', 'Sostenibilidad, ecología y cuidado del planeta', 6, 6, NOW(), NOW()),
                                                                                           ('Hablemos de Psicología 🧠', 'Bienestar mental, desarrollo personal y salud emocional', 7, 7, NOW(), NOW()),
                                                                                           ('Historia 📚', 'Explora el pasado y descubre curiosidades históricas', 8, 8, NOW(), NOW()),
                                                                                           ('Astronomía 🌌', 'Explora el universo, estrellas y fenómenos astronómicos', 9, 9, NOW(), NOW()),
                                                                                           ('Filosofía 🤔', 'Reflexiones filosóficas, ética y pensamiento crítico', 10, 10, NOW(), NOW()),
                                                                                           ('Cueva gamer 🕹️', 'Club para jugar videojuegos y compartir experiencias gaming', 11, 11, NOW(), NOW()),
                                                                                           ('Club de lectura 📚', 'Espacio para compartir reseñas de libros, recomendaciones literarias y discutir sobre nuestras lecturas favoritas', 12, 12, NOW(), NOW()),
                                                                                           ('General UNQUI', 'Espacio para hablar de cualquier cosa sobre la universidad', 13, 13, NOW(), NOW()),
                                                                                           ('Tecnología 💻', 'Discusiones sobre las últimas tendencias tecnológicas y programación', 1, 1, NOW(), NOW()),
                                                                                           ('Club de Música 🎵', 'Comparte tus canciones favoritas y descubre nueva música', 2, 2, NOW(), NOW()),
                                                                                           ('Deportes ⚽', 'Todo sobre deportes, desde fútbol hasta tenis', 3, 3, NOW(), NOW()),
                                                                                           ('Arte y Diseño 🎨', 'Comparte tus creaciones artísticas y proyectos de diseño', 4, 4, NOW(), NOW()),
                                                                                           ('Grupo de Cocina 👨‍🍳', 'Recetas, técnicas culinarias y experiencias gastronómicas', 5, 5, NOW(), NOW()),
                                                                                           ('Club de viajeros ✈️', 'Comparte tus experiencias de viaje y descubre nuevos destinos', 6, 6, NOW(), NOW()),
                                                                                           ('Solo para cientificos 🔬', 'Discusiones sobre avances científicos y curiosidades del mundo', 7, 7, NOW(), NOW());



-- Insert sample posts
INSERT INTO posts (title, content, created_by, space_id, updated_by, created_at, updated_at) VALUES
                                                                                                 ('¿Cuál es el mejor framework para React en 2024?', 'Estoy evaluando opciones para un nuevo proyecto. ¿Qué opinan sobre Next.js vs Vite vs Create React App?', 1, 1, 1, NOW(), NOW()),
                                                                                                 ('Configuración de Docker para microservicios', 'Comparto mi experiencia configurando Docker Compose para un proyecto con múltiples servicios', 2, 2, 2, NOW(), NOW()),
                                                                                                 ('Introducción a TensorFlow 2.0', 'Guía básica para comenzar con machine learning usando TensorFlow', 3, 3, 3, NOW(), NOW()),
                                                                                                 ('Flutter vs React Native: Comparación 2024', 'Análisis detallado de ambas tecnologías para desarrollo móvil multiplataforma', 4, 4, 4, NOW(), NOW()),
                                                                                                 ('Arquitectura de microservicios con Go', 'Cómo diseñar una arquitectura escalable usando Go y gRPC', 5, 5, 5, NOW(), NOW()),
                                                                                                 ('Optimización de rendimiento en React', 'Técnicas avanzadas para mejorar el rendimiento de aplicaciones React', 6, 6, 6, NOW(), NOW());

-- Insert sample comments
INSERT INTO comments (post_id, content, created_by, created_at) VALUES
                                                                    (1, 'Yo recomendaría Next.js para proyectos que necesiten SSR', 2, NOW()),
                                                                    (1, 'Vite es excelente para desarrollo, muy rápido', 6, NOW()),
                                                                    (2, 'Muy útil, gracias por compartir la configuración', 5, NOW()),
                                                                    (2, '¿Has probado con Kubernetes en lugar de Docker Compose?', 7, NOW()),
                                                                    (3, 'Excelente guía, muy clara para principiantes', 8, NOW()),
                                                                    (3, '¿Podrías agregar ejemplos de redes neuronales?', 1, NOW()),
                                                                    (4, 'React Native tiene mejor ecosistema, pero Flutter es más consistente', 5, NOW()),
                                                                    (4, 'Flutter tiene mejor rendimiento nativo', 7, NOW()),
                                                                    (5, 'Go es perfecto para microservicios, muy eficiente', 3, NOW()),
                                                                    (5, '¿Usas algún ORM específico con Go?', 8, NOW()),
                                                                    (6, 'React.memo es clave para optimización', 1, NOW()),
                                                                    (6, 'También recomiendo usar React DevTools Profiler', 4, NOW());

-- =========================
-- POSTS 1–20
-- =========================
INSERT INTO posts (title, content, created_by, space_id, updated_by, created_at, updated_at) VALUES
-- Fotografía (espacio 1)
('Mejores cámaras para empezar en fotografía', '¿Qué cámara recomiendan para principiantes en 2024?', 1, 1, 1, NOW(), NOW()),
('Trucos de edición en Lightroom', 'Comparto algunos atajos y técnicas para editar más rápido', 5, 1, 5, NOW(), NOW()),

-- Fitness (espacio 2)
('Rutina full body para principiantes', 'Una guía para entrenar todo el cuerpo tres veces por semana', 2, 2, 2, NOW(), NOW()),
('Alimentos clave para ganar masa muscular', 'Comparto una lista de comidas que me ayudaron a subir de peso sano', 8, 2, 8, NOW(), NOW()),

-- Cine y series (espacio 3)
('Mejores películas de ciencia ficción', '¿Cuáles son sus top 5 de sci-fi?', 3, 3, 3, NOW(), NOW()),
('Series cortas pero intensas', 'Busco recomendaciones de miniseries de menos de 10 capítulos', 6, 3, 6, NOW(), NOW()),

-- Emprendimientos (espacio 4)
('¿Vale la pena emprender solo?', 'Estoy evaluando lanzar un proyecto sin socios. ¿Qué opinan?', 4, 4, 4, NOW(), NOW()),
('Apps útiles para startups', 'Compilo una lista de herramientas que me ayudaron en la gestión diaria', 10, 4, 10, NOW(), NOW()),

-- Idiomas (espacio 5)
('Cómo mejorar la pronunciación en inglés', 'Tips prácticos para sonar más natural', 7, 5, 7, NOW(), NOW()),
('Apps para aprender japonés', '¿Alguien probó Duolingo, LingQ o WaniKani? ¿Cuál recomiendan?', 13, 5, 13, NOW(), NOW()),

-- Medio ambiente (espacio 6)
('Huerta urbana en balcones pequeños', 'Comparto mi experiencia cultivando tomates en un mini balcón', 9, 6, 9, NOW(), NOW()),
('Consejos para reducir plásticos en casa', 'Ideas simples que cualquiera puede aplicar en su día a día', 11, 6, 11, NOW(), NOW()),

-- Psicología (espacio 7)
('Cómo gestionar la ansiedad', 'Estrategias que me funcionaron en momentos de estrés', 12, 7, 12, NOW(), NOW()),
('Libros recomendados de psicología', 'Estoy buscando lecturas de divulgación accesibles', 2, 7, 2, NOW(), NOW()),

-- Historia (espacio 8)
('Curiosidades de la Edad Media', 'Pequeños datos históricos sorprendentes que descubrí leyendo', 14, 8, 14, NOW(), NOW()),
('Revolución Industrial: pros y contras', 'Discusión sobre su impacto en la sociedad', 5, 8, 5, NOW(), NOW()),

-- Astronomía (espacio 9)
('Eclipse solar 2024', '¿Quién planea viajar a verlo en directo?', 3, 9, 3, NOW(), NOW()),
('Telescopios para principiantes', 'Recomendaciones para iniciarse en la observación astronómica', 7, 9, 7, NOW(), NOW()),

-- Filosofía (espacio 10)
('¿Existe el libre albedrío?', 'Un debate clásico que sigue abierto', 6, 10, 6, NOW(), NOW()),
('Filosofía estoica en la vida moderna', 'Cómo aplico el estoicismo en mi rutina diaria', 1, 10, 1, NOW(), NOW()),

-- Cueva gamer (espacio 11)
('Mejores juegos cooperativos online', 'Busco recomendaciones para jugar con amigos en PC', 8, 11, 8, NOW(), NOW()),
('¿Vale la pena la PS5 en 2024?', 'Opiniones sobre si ya conviene o esperar la siguiente gen', 10, 11, 10, NOW(), NOW());

-- =========================
-- COMMENTS (cada post tiene 2–3)
-- =========================
INSERT INTO comments (post_id, content, created_by, created_at) VALUES
                                                                    (1, 'Para empezar recomiendo Canon EOS Rebel', 2, NOW()),
                                                                    (1, 'Sony A6000 es una opción económica', 6, NOW()),
                                                                    (2, '¡Los atajos de teclado salvan horas!', 3, NOW()),
                                                                    (3, 'Rutina sencilla y efectiva, gracias', 5, NOW()),
                                                                    (3, '¿Cuántas series recomiendas por ejercicio?', 7, NOW()),
                                                                    (4, 'El arroz con pollo es clave en volumen', 9, NOW()),
                                                                    (5, 'Mi top incluye Interestelar y Blade Runner', 4, NOW()),
                                                                    (5, '¿Alguien vio Dune 2? Brutal', 11, NOW()),
                                                                    (6, 'Chernobyl es la mejor miniserie que vi', 13, NOW()),
                                                                    (6, 'También recomiendo The Night Of', 14, NOW()),
                                                                    (7, 'Emprender solo es más arriesgado', 12, NOW()),
                                                                    (7, 'Depende de tu tolerancia al riesgo', 8, NOW()),
                                                                    (8, 'Notion y Trello me salvaron la vida', 1, NOW()),
                                                                    (9, 'Hablar en voz alta ayuda mucho', 10, NOW()),
                                                                    (10, 'WaniKani es excelente para kanjis', 5, NOW()),
                                                                    (10, 'Duolingo sirve pero se queda corto', 2, NOW()),
                                                                    (11, 'Yo cultivo albahaca en macetas pequeñas', 4, NOW()),
                                                                    (12, 'Cambiar botellas plásticas por termo metálico', 6, NOW()),
                                                                    (13, 'La meditación me ayudó bastante', 7, NOW()),
                                                                    (14, 'Recomiendo “El cerebro y la inteligencia emocional”', 9, NOW()),
                                                                    (15, 'Los castillos medievales eran más pequeños de lo que pensamos', 11, NOW()),
                                                                    (16, 'Aceleró la urbanización, pero explotó a los trabajadores', 2, NOW()),
                                                                    (17, '¡Yo viajo a México para verlo!', 14, NOW()),
                                                                    (18, 'SkyWatcher 130 es muy bueno calidad/precio', 3, NOW()),
                                                                    (19, 'El libre albedrío es una ilusión', 10, NOW()),
                                                                    (19, 'No estoy de acuerdo, siempre tenemos elección', 13, NOW()),
                                                                    (20, 'Marco Aurelio es una gran referencia', 12, NOW()),
                                                                    (21, 'It Takes Two es imperdible', 9, NOW()),
                                                                    (21, 'También recomiendo Deep Rock Galactic', 11, NOW()),
                                                                    (22, 'Yo la compré y no me arrepiento', 6, NOW()),
                                                                    (22, 'Mejor esperar una revisión Slim', 4, NOW());


-- =========================
-- POSTS 21–40
-- =========================
INSERT INTO posts (title, content, created_by, space_id, updated_by, created_at, updated_at) VALUES
-- Fotografía (espacio 1)
('Fotografía nocturna: tips básicos', 'Cómo mejorar fotos de la luna y estrellas con cámaras básicas', 4, 1, 4, NOW(), NOW()),
('Diferencias entre RAW y JPG', '¿Conviene disparar siempre en RAW?', 12, 1, 12, NOW(), NOW()),

-- Fitness (espacio 2)
('Cardio en ayunas: ¿mito o realidad?', 'He leído opiniones encontradas, ¿qué piensan?', 6, 2, 6, NOW(), NOW()),
('Suplementos recomendados para principiantes', 'Proteínas, creatina, pre entreno... ¿cuáles valen la pena?', 11, 2, 11, NOW(), NOW()),

-- Cine y series (espacio 3)
('¿Cuál es la mejor película de Christopher Nolan?', 'Me cuesta decidir entre Inception, Interstellar y Oppenheimer', 8, 3, 8, NOW(), NOW()),
('Series infravaloradas en Netflix', 'Comparto algunas joyitas poco conocidas', 10, 3, 10, NOW(), NOW()),

-- Emprendimientos (espacio 4)
('Cómo conseguir inversores', 'Consejos para presentar tu pitch a fondos de inversión', 5, 4, 5, NOW(), NOW()),
('Errores comunes al emprender', 'Lo que aprendí después de dos startups fallidas', 13, 4, 13, NOW(), NOW()),

-- Idiomas (espacio 5)
('¿Vale la pena un profesor particular?', 'Comparando apps de idiomas vs clases personalizadas', 9, 5, 9, NOW(), NOW()),
('Mejores podcasts para practicar francés', 'Busco contenido para escuchar en el auto', 2, 5, 2, NOW(), NOW()),

-- Medio ambiente (espacio 6)
('Paneles solares en casa: mi experiencia', 'Instalé un kit básico y lo cuento aquí', 7, 6, 7, NOW(), NOW()),
('¿El auto eléctrico es realmente ecológico?', 'Debate sobre la huella de producción de baterías', 3, 6, 3, NOW(), NOW()),

-- Psicología (espacio 7)
('Mindfulness explicado fácil', 'Una práctica que me ayudó a enfocarme en el presente', 14, 7, 14, NOW(), NOW()),
('Cómo ayudar a un amigo con depresión', 'Busco consejos prácticos para acompañar', 1, 7, 1, NOW(), NOW()),

-- Historia (espacio 8)
('Roma y su legado en el derecho', 'Cómo influyó en los sistemas legales actuales', 10, 8, 10, NOW(), NOW()),
('Segunda Guerra Mundial: libros recomendados', 'Quiero leer más allá de los clásicos', 12, 8, 12, NOW(), NOW()),

-- Astronomía (espacio 9)
('Las lunas de Júpiter', 'Datos curiosos sobre Io, Europa, Ganímedes y Calisto', 11, 9, 11, NOW(), NOW()),
('La paradoja de Fermi', 'Si el universo es tan grande, ¿dónde están todos?', 9, 9, 9, NOW(), NOW()),

-- Filosofía (espacio 10)
('Nietzsche y el superhombre', 'Cómo interpretar este concepto hoy en día', 5, 10, 5, NOW(), NOW()),
('Ética de la inteligencia artificial', 'Reflexiones sobre los límites morales del uso de IA', 8, 10, 8, NOW(), NOW()),

-- Cueva gamer (espacio 11)
('Mejores RPGs de la última década', '¿Cuál creen que es el top 3?', 2, 11, 2, NOW(), NOW()),
('El futuro del gaming en la nube', '¿Conviene apostar por GeForce Now, Xbox Cloud o PS Plus?', 6, 11, 6, NOW(), NOW());

-- =========================
-- COMMENTS 21–40
-- =========================
INSERT INTO comments (post_id, content, created_by, created_at) VALUES
                                                                    (23, 'Un trípode ayuda muchísimo en nocturnas', 7, NOW()),
                                                                    (23, 'También usar ISO bajo y exposición larga', 9, NOW()),
                                                                    (24, 'RAW da más control, pero pesa mucho', 3, NOW()),
                                                                    (24, 'JPG es suficiente para redes sociales', 13, NOW()),
                                                                    (25, 'El cardio en ayunas no es mágico, solo un método más', 8, NOW()),
                                                                    (25, 'Yo lo probé y me ayudó a bajar grasa', 12, NOW()),
                                                                    (26, 'La creatina es lo más respaldado por estudios', 1, NOW()),
                                                                    (26, 'Proteína en polvo solo si no llegas con comida real', 14, NOW()),
                                                                    (27, 'Inception me voló la cabeza', 4, NOW()),
                                                                    (27, 'Interstellar tiene mejor banda sonora', 11, NOW()),
                                                                    (28, 'Dark es de lo mejor que vi en Netflix', 10, NOW()),
                                                                    (28, 'The OA está muy infravalorada', 2, NOW()),
                                                                    (29, 'Un pitch claro y breve es clave', 6, NOW()),
                                                                    (29, 'Conseguir tracción antes ayuda mucho', 7, NOW()),
                                                                    (30, 'Fracasar enseña más que el éxito', 9, NOW()),
                                                                    (30, 'Totalmente, cada error es aprendizaje', 5, NOW()),
                                                                    (31, 'Un profe corrige errores que una app no', 11, NOW()),
                                                                    (31, 'Depende del presupuesto y disciplina', 3, NOW()),
                                                                    (32, 'Coffee Break French es muy bueno', 12, NOW()),
                                                                    (32, 'También InnerFrench', 14, NOW()),
                                                                    (33, '¿Cuánto tardaste en recuperar la inversión?', 8, NOW()),
                                                                    (33, 'Quiero poner paneles pero me da miedo el costo', 2, NOW()),
                                                                    (34, 'La minería de litio es un gran problema', 4, NOW()),
                                                                    (34, 'El balance aún no es tan verde como parece', 10, NOW()),
                                                                    (35, 'Mindfulness me ayudó a dormir mejor', 7, NOW()),
                                                                    (36, 'Lo más importante es escuchar sin juzgar', 13, NOW()),
                                                                    (36, 'A veces solo acompañar ya es suficiente', 9, NOW()),
                                                                    (37, 'El derecho romano está en nuestras leyes civiles', 1, NOW()),
                                                                    (38, 'Recomiendo “El Tercer Reich” de Shirer', 6, NOW()),
                                                                    (39, 'Europa es candidata para albergar vida', 12, NOW()),
                                                                    (39, 'Ganímedes es enorme, casi un planeta', 5, NOW()),
                                                                    (40, 'La paradoja es fascinante, tal vez seamos los primeros', 8, NOW()),
                                                                    (41, 'Nietzsche inspira pero también se malinterpreta mucho', 4, NOW()),
                                                                    (42, 'La IA debe usarse con responsabilidad', 2, NOW()),
                                                                    (42, 'El problema es que la ética avanza más lento que la tecnología', 10, NOW()),
                                                                    (43, 'The Witcher 3 está en mi top', 9, NOW()),
                                                                    (43, 'Persona 5 es otra joya', 13, NOW()),
                                                                    (44, 'Xbox Cloud tiene mejor estabilidad', 14, NOW()),
                                                                    (44, 'GeForce Now va bien si tienes buena conexión', 11, NOW());


-- =========================
-- POSTS 41–70
-- =========================
INSERT INTO posts (title, content, created_by, space_id, updated_by, created_at, updated_at) VALUES
-- Fotografía (espacio 1)
('Fotografía de retrato: iluminación natural vs artificial', 'Ventajas y desventajas de cada tipo de luz', 3, 1, 3, NOW(), NOW()),
('Cómo elegir el lente adecuado', 'Guía rápida para diferentes tipos de fotografía', 7, 1, 7, NOW(), NOW()),

-- Fitness (espacio 2)
('Entrenamiento HIIT en casa', 'Rutina rápida de 20 minutos para quemar grasa', 4, 2, 4, NOW(), NOW()),
('Estiramientos para después del entrenamiento', 'Evitar lesiones y mejorar la recuperación', 12, 2, 12, NOW(), NOW()),

-- Cine y series (espacio 3)
('Películas de terror que valen la pena', 'Busco recomendaciones más allá de los clásicos', 1, 3, 1, NOW(), NOW()),
('Documentales imperdibles', 'Desde naturaleza hasta historia, ¿cuáles recomiendan?', 5, 3, 5, NOW(), NOW()),

-- Emprendimientos (espacio 4)
('Marketing digital para startups', 'Herramientas y estrategias que funcionan', 6, 4, 6, NOW(), NOW()),
('Cómo validar tu idea antes de invertir', 'Evitar gastar dinero en productos que nadie quiere', 2, 4, 2, NOW(), NOW()),

-- Idiomas (espacio 5)
('Intercambio de idiomas online', 'Plataformas y tips para practicar con nativos', 9, 5, 9, NOW(), NOW()),
('Gramática inglesa: trucos para no fallar', 'Errores comunes que todo el mundo comete', 14, 5, 14, NOW(), NOW()),

-- Medio ambiente (espacio 6)
('Bicicleta vs auto: impacto ambiental', 'Comparando emisiones y beneficios', 3, 6, 3, NOW(), NOW()),
('Reciclaje de electrónicos', 'Qué se puede reciclar y cómo', 5, 6, 5, NOW(), NOW()),

-- Psicología (espacio 7)
('Técnicas para mejorar la memoria', 'Ejercicios simples para entrenar el cerebro', 6, 7, 6, NOW(), NOW()),
('Cómo manejar el estrés laboral', 'Estrategias prácticas para el día a día', 8, 7, 8, NOW(), NOW()),

-- Historia (espacio 8)
('Imperio Otomano: datos curiosos', 'Cultura, política y costumbres que sorprenden', 10, 8, 10, NOW(), NOW()),
('Revoluciones americanas vs francesas', 'Comparando causas y consecuencias', 12, 8, 12, NOW(), NOW()),

-- Astronomía (espacio 9)
('Cometas famosos en la historia', 'Halley, Hale-Bopp y otros', 1, 9, 1, NOW(), NOW()),
('Constelaciones fáciles de identificar', 'Tips para principiantes', 4, 9, 4, NOW(), NOW()),

-- Filosofía (espacio 10)
('El existencialismo en la literatura', 'Autores y libros recomendados', 7, 10, 7, NOW(), NOW()),
('Ética ambiental y responsabilidad', 'Reflexión sobre nuestro impacto en la naturaleza', 13, 10, 13, NOW(), NOW()),

-- Cueva gamer (espacio 11)
('Juegos indie que merecen atención', 'Recomendaciones de títulos menos conocidos', 2, 11, 2, NOW(), NOW()),
('Hardware gaming económico', 'PC o consolas para presupuestos bajos', 11, 11, 11, NOW(), NOW()),
('Streaming de partidas: consejos', 'Cómo empezar a transmitir sin gastar mucho', 6, 11, 6, NOW(), NOW());

-- =========================
-- COMMENTS 41–70
-- =========================
INSERT INTO comments (post_id, content, created_by, created_at) VALUES
                                                                    (45, 'La luz natural siempre da un look más suave', 2, NOW()),
                                                                    (45, 'Yo prefiero flash para retratos dramáticos', 9, NOW()),
                                                                    (46, 'Depende si quieres paisaje o retrato', 3, NOW()),
                                                                    (46, 'Un 50mm f1.8 es versátil y barato', 14, NOW()),
                                                                    (47, 'HIIT es brutal si tienes poco tiempo', 1, NOW()),
                                                                    (47, 'Cuidado con la técnica, evita lesiones', 7, NOW()),
                                                                    (48, 'Estirar después es clave', 12, NOW()),
                                                                    (48, 'También ayuda hacer foam roller', 5, NOW()),
                                                                    (49, 'It sigue siendo mi favorita del año', 6, NOW()),
                                                                    (49, 'La bruja de Blair es un clásico eterno', 11, NOW()),
                                                                    (50, 'Planet Earth II es imperdible', 8, NOW()),
                                                                    (50, 'Recomiendo “The Social Dilemma” para tech lovers', 13, NOW()),
                                                                    (51, 'Google Analytics y Mailchimp son útiles', 9, NOW()),
                                                                    (51, 'También recomiendo Canva para marketing rápido', 2, NOW()),
                                                                    (52, 'Hacer encuestas y prototipos ayuda mucho', 1, NOW()),
                                                                    (52, 'Nunca subestimes el feedback de usuarios reales', 12, NOW()),
                                                                    (53, 'HelloTalk es muy bueno para practicar idiomas', 7, NOW()),
                                                                    (53, 'Tandem también es útil', 6, NOW()),
                                                                    (54, 'Los errores con “their/there/they’re” son comunes', 14, NOW()),
                                                                    (54, 'No olvidar el uso de tiempos verbales', 3, NOW()),
                                                                    (55, 'Bicicleta siempre que puedas', 8, NOW()),
                                                                    (55, 'El auto eléctrico reduce emisiones pero no todo es verde', 5, NOW()),
                                                                    (56, 'Llevar baterías viejas a puntos limpios', 9, NOW()),
                                                                    (56, 'No botar celulares a la basura', 4, NOW()),
                                                                    (57, 'Ejercicios de memoria con cartas funcionan', 10, NOW()),
                                                                    (57, 'Repetir nombres y listas ayuda mucho', 1, NOW()),
                                                                    (58, 'Respirar profundo antes de empezar', 12, NOW()),
                                                                    (58, 'Organizar tareas y pausas también ayuda', 2, NOW()),
                                                                    (59, 'La arquitectura otomana es impresionante', 6, NOW()),
                                                                    (59, 'Sus bazares eran únicos', 11, NOW()),
                                                                    (60, 'Ambas revoluciones cambiaron el mundo', 14, NOW()),
                                                                    (60, 'La revolución francesa fue más radical', 3, NOW()),
                                                                    (61, 'Halley aparece cada 76 años', 7, NOW()),
                                                                    (61, 'Hale-Bopp dejó un espectáculo inolvidable', 10, NOW()),
                                                                    (62, 'Orión es fácil de encontrar', 2, NOW()),
                                                                    (62, 'También busca Casiopea en el cielo norte', 12, NOW()),
                                                                    (63, 'Camus y Sartre son imperdibles', 5, NOW()),
                                                                    (63, '“El extranjero” es excelente para empezar', 8, NOW()),
                                                                    (64, 'La ética ambiental debe ser prioritaria', 6, NOW()),
                                                                    (64, 'Todos podemos hacer cambios pequeños diarios', 9, NOW()),
                                                                    (65, 'Hollow Knight es un juego indie genial', 1, NOW()),
                                                                    (66, 'Una GTX 1650 sirve para empezar', 4, NOW()),
                                                                    (67, 'OBS Studio es gratis y potente', 11, NOW()),
                                                                    (67, 'Recomiendo aprender a usar overlays', 13, NOW());


