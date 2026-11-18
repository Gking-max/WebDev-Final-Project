INSERT INTO statuses (name) VALUES ('New'),('Assigned'),('In progress'),('Waiting for you'),('Resolved') ON CONFLICT (name) DO NOTHING;
INSERT INTO priorities (name) VALUES ('Low'),('Medium'),('High'),('Urgent') ON CONFLICT (name) DO NOTHING;
INSERT INTO categories (name) VALUES ('Network'),('Email'),('Hardware'),('Software'),('Account access') ON CONFLICT (name) DO NOTHING;

INSERT INTO tickets (title, description, requester_name, requester_type, location_name, location_code, status_id, priority_id, category_id)
SELECT * FROM (VALUES
  ('Wi-Fi not working in Library','Cannot connect to campus Wi-Fi from the second floor.','Alex Smith','Student','Main Library','LIB-001',(SELECT id FROM statuses WHERE name='New'),(SELECT id FROM priorities WHERE name='High'),(SELECT id FROM categories WHERE name='Network')),
  ('Cannot access email','Authentication error when signing in.','Janelle Brown','Student','Student Center','STU-201',(SELECT id FROM statuses WHERE name='Assigned'),(SELECT id FROM priorities WHERE name='Medium'),(SELECT id FROM categories WHERE name='Email')),
  ('Printer not printing','Paper jam error on display.','Marcus Lee','Faculty','ICT Building','ICT-105',(SELECT id FROM statuses WHERE name='In progress'),(SELECT id FROM priorities WHERE name='Medium'),(SELECT id FROM categories WHERE name='Hardware')),
  ('Laptop running very slowly after the latest operating system update','Performance issues, slow boot.','Sofia Grant','Student','Science Building','SCI-203',(SELECT id FROM statuses WHERE name='Waiting for you'),(SELECT id FROM priorities WHERE name='High'),(SELECT id FROM categories WHERE name='Software')),
  ('Request for software license','Need MATLAB access.','Daniel Roberts','Student','Online / Remote','REMOTE',(SELECT id FROM statuses WHERE name='Resolved'),(SELECT id FROM priorities WHERE name='Low'),(SELECT id FROM categories WHERE name='Software')),
  ('Password reset for staff portal','Locked out after failed attempts.','Priya Nair','Staff','Administration','ADM-012',(SELECT id FROM statuses WHERE name='New'),(SELECT id FROM priorities WHERE name='Urgent'),(SELECT id FROM categories WHERE name='Account access')),
  ('Projector flickering in lecture hall','Flickers during presentations.','Tomás Reyes','Faculty','Engineering Block','ENG-204',(SELECT id FROM statuses WHERE name='In progress'),(SELECT id FROM priorities WHERE name='Medium'),(SELECT id FROM categories WHERE name='Hardware')),
  ('Cannot log into the learning management system from off campus','Off-campus login times out.','Hana Okafor','Student','Online / Remote','REMOTE',(SELECT id FROM statuses WHERE name='Assigned'),(SELECT id FROM priorities WHERE name='High'),(SELECT id FROM categories WHERE name='Network'))
) AS v
WHERE NOT EXISTS (SELECT 1 FROM tickets);