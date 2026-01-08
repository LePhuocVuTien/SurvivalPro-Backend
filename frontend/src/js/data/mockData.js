// Dữ liệu mẫu cho ứng dụng

window.mockData = {
    users: [
        { id: 1, email: 'user1@example.com', password: 'hashed_password_1', name: 'John Doe', avatar: 'https://randomuser.me/api/portraits/men/32.jpg', push_token: 'push_token_1', created_at: '2023-10-15T08:30:00Z' },
        { id: 2, email: 'user2@example.com', password: 'hashed_password_2', name: 'Jane Smith', avatar: 'https://randomuser.me/api/portraits/women/44.jpg', push_token: 'push_token_2', created_at: '2023-10-16T09:15:00Z' },
        { id: 3, email: 'user3@example.com', password: 'hashed_password_3', name: 'Robert Johnson', avatar: 'https://randomuser.me/api/portraits/men/67.jpg', push_token: 'push_token_3', created_at: '2023-10-17T10:45:00Z' },
        { id: 4, email: 'user4@example.com', password: 'hashed_password_4', name: 'Emily Davis', avatar: 'https://randomuser.me/api/portraits/women/65.jpg', push_token: null, created_at: '2023-10-18T14:20:00Z' },
        { id: 5, email: 'user5@example.com', password: 'hashed_password_5', name: 'Michael Wilson', avatar: 'https://randomuser.me/api/portraits/men/31.jpg', push_token: 'push_token_5', created_at: '2023-10-19T16:10:00Z' }
    ],
    
    checklist: [
        { id: 1, user_id: 1, title: 'First Aid Kit', category: 'supplies', description: 'Complete first aid kit with bandages, antiseptics, and medications', is_checked: true, created_at: '2023-10-20T08:30:00Z' },
        { id: 2, user_id: 1, title: 'Emergency Documents', category: 'documents', description: 'Passport, ID cards, insurance papers, and emergency contacts', is_checked: false, created_at: '2023-10-20T09:15:00Z' },
        { id: 3, user_id: 2, title: 'Water Supply', category: 'water', description: 'At least 3 gallons of water per person', is_checked: true, created_at: '2023-10-21T10:30:00Z' },
        { id: 4, user_id: 2, title: 'Non-perishable Food', category: 'food', description: 'Canned goods, energy bars, and dried fruits', is_checked: true, created_at: '2023-10-21T11:45:00Z' },
        { id: 5, user_id: 3, title: 'Emergency Radio', category: 'emergency', description: 'Hand-crank radio with weather alerts', is_checked: false, created_at: '2023-10-22T13:20:00Z' },
        { id: 6, user_id: 4, title: 'Tent and Sleeping Bags', category: 'shelter', description: 'Emergency shelter for 4 people', is_checked: true, created_at: '2023-10-23T15:10:00Z' },
        { id: 7, user_id: 5, title: 'Flashlight and Batteries', category: 'supplies', description: 'LED flashlight with extra batteries', is_checked: false, created_at: '2023-10-24T16:40:00Z' }
    ],
    
    locations: [
        { id: 1, user_id: 1, latitude: 21.028511, longtitude: 105.804817, created_at: '2023-10-25T08:15:00Z' }, // Hanoi
        { id: 2, user_id: 2, latitude: 10.823099, longtitude: 106.629662, created_at: '2023-10-25T09:30:00Z' }, // Ho Chi Minh City
        { id: 3, user_id: 3, latitude: 16.054407, longtitude: 108.202164, created_at: '2023-10-25T10:45:00Z' }, // Da Nang
        { id: 4, user_id: 4, latitude: 20.844912, longtitude: 106.688084, created_at: '2023-10-25T11:20:00Z' }, // Hai Phong
        { id: 5, user_id: 5, latitude: 12.238791, longtitude: 109.196749, created_at: '2023-10-25T12:10:00Z' }, // Nha Trang
        { id: 6, user_id: 1, latitude: 21.033333, longtitude: 105.849998, created_at: '2023-10-26T08:45:00Z' }, // Hanoi - updated
        { id: 7, user_id: 2, latitude: 10.776889, longtitude: 106.700806, created_at: '2023-10-26T10:30:00Z' }  // Ho Chi Minh City - updated
    ],

    guides: [
        {
            id: 1,
            title: 'Sơ cứu vết thương cơ bản',
            category: 'first-aid',
            difficulty: 'easy',
            icon: '🩹',
            content: 'Các bước sơ cứu vết thương:\n\n1. Rửa tay sạch trước khi xử lý\n2. Làm sạch vết thương bằng nước sạch\n3. Bôi thuốc sát trùng (betadine, cồn)\n4. Băng gạc vô trùng\n5. Thay băng hàng ngày\n\nLưu ý: Nếu vết thương sâu hoặc chảy máu nhiều, cần đến cơ sở y tế ngay lập tức.',
            image_url: 'https://images.unsplash.com/photo-1603398938378-e54eab446dde?w=400',
            views: 1250,
            created_at: '2024-01-15T10:30:00Z'
        },
        {
            id: 2,
            title: 'Cách dựng lều tạm bợ trong rừng',
            category: 'shelter',
            difficulty: 'medium',
            icon: '⛺',
            content: 'Hướng dẫn dựng lều khẩn cấp:\n\n1. Chọn địa điểm: Tìm nơi cao ráo, tránh đáy thung lũng\n2. Tìm vật liệu: Cành cây, lá to, dây leo\n3. Dựng khung: Tạo khung hình chữ A hoặc dựa vào cây\n4. Phủ lớp che: Dùng lá cây xếp từ dưới lên trên\n5. Gia cố: Dùng dây buộc chắc chắn\n\nMẹo: Hướng cửa lều tránh gió, tạo rãnh thoát nước xung quanh.',
            image_url: 'https://images.unsplash.com/photo-1504280390367-361c6d9f38f4?w=400',
            views: 890,
            created_at: '2024-01-15T11:00:00Z'
        },
        {
            id: 3,
            title: 'Tìm và lọc nước uống an toàn',
            category: 'food',
            difficulty: 'medium',
            icon: '💧',
            content: 'Các phương pháp lọc nước:\n\n1. Tìm nguồn: Ưu tiên nước chảy, suối, sông\n2. Lọc sơ bộ: Dùng vải lọc bỏ cặn bẩn\n3. Đun sôi: Đun sôi ít nhất 5-10 phút\n4. Lọc than: Dùng than củi nghiền nhỏ để khử mùi\n5. Phơi nắng: Nếu có, phơi nắng UV 6 giờ\n\nDấu hiệu nước an toàn: Trong, không mùi lạ, không vị lạ.',
            image_url: 'https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=400',
            views: 2150,
            created_at: '2024-01-16T08:30:00Z'
        },
        {
            id: 4,
            title: 'Định hướng bằng la bàn và bản đồ',
            category: 'navigation',
            difficulty: 'medium',
            icon: '🧭',
            content: 'Kỹ năng định hướng cơ bản:\n\n1. Đọc bản đồ: Hiểu ký hiệu, tỷ lệ, đường đồng mức\n2. Sử dụng la bàn: Xác định hướng Bắc từ\n3. Định vị: Dùng 2-3 điểm mốc để tam giác hóa\n4. Đi theo azimuth: Giữ hướng cố định\n5. Nhận dạng địa hình: So sánh thực tế với bản đồ\n\nKhông có la bàn: Dùng mặt trời, ngôi sao, rêu cây.',
            image_url: 'https://images.unsplash.com/photo-1569163139394-de4798aa62b6?w=400',
            views: 650,
            created_at: '2024-01-16T14:15:00Z'
        },
        {
            id: 5,
            title: 'Nhóm lửa trong điều kiện khắc nghiệt',
            category: 'fire',
            difficulty: 'hard',
            icon: '🔥',
            content: 'Cách nhóm lửa khi khó khăn:\n\n1. Chuẩn bị: Củi khô, phoi bào, cành nhỏ\n2. Tạo tinder: Vỏ cây khô, lông tơ, giấy\n3. Xếp củi: Kiểu chữ A hoặc hình tháp\n4. Nhóm lửa: Dùng diêm, lửa cọ, thấu kính\n5. Duy trì: Thêm củi từ nhỏ đến lớn\n\nDưới mưa: Tìm củi khô bên trong thân cây gãy, dùng nhựa cây.',
            image_url: 'https://images.unsplash.com/photo-1525498128493-380d1990a112?w=400',
            views: 1450,
            created_at: '2024-01-17T09:00:00Z'
        },
        {
            id: 6,
            title: 'Tín hiệu cấp cứu SOS',
            category: 'signaling',
            difficulty: 'easy',
            icon: '🆘',
            content: 'Các phương pháp báo hiệu:\n\n1. Tín hiệu âm thanh:\n   - SOS: 3 tiếng ngắn, 3 tiếng dài, 3 tiếng ngắn\n   - Còi, huýt sáo mỗi 10 giây\n\n2. Tín hiệu ánh sáng:\n   - Gương phản chiếu\n   - Đèn pin nhấp nháy theo SOS\n\n3. Tín hiệu khói:\n   - Ban ngày: Khói đen (cao su, lá xanh)\n   - Ban đêm: Lửa sáng\n\n4. Tín hiệu mặt đất:\n   - Chữ X: Cần cứu trợ\n   - Chữ V: Cần hỗ trợ y tế\n   - Tam giác: An toàn',
            image_url: 'https://images.unsplash.com/photo-1584438784894-089d6a62b8fa?w=400',
            views: 980,
            created_at: '2024-01-17T11:30:00Z'
        }
    ],
    
    notifications: [
        { id: 1, user_id: 1, title: 'Checklist Reminder', body: 'Don\'t forget to check your emergency supplies this week', data: '{"type": "reminder", "checklist_id": 1}', type: 'push', is_read: true, sent: true, sent_at: '2023-10-25T09:00:00Z', created_at: '2023-10-24T14:30:00Z' },
        { id: 2, user_id: 2, title: 'New Survival Guide', body: 'A new guide "Water Purification Methods" is now available', data: '{"type": "new_guide", "guide_id": 3}', type: 'in_app', is_read: false, sent: true, sent_at: '2023-10-25T10:15:00Z', created_at: '2023-10-24T16:45:00Z' },
        { id: 3, user_id: 3, title: 'Weather Alert', body: 'Severe weather warning in your area. Please take precautions.', data: '{"type": "weather_alert", "severity": "high"}', type: 'both', is_read: true, sent: true, sent_at: '2023-10-26T07:30:00Z', created_at: '2023-10-25T18:20:00Z' },
        { id: 4, user_id: 4, title: 'Location Sharing', body: 'Your location is being shared with emergency contacts', data: '{"type": "location_share", "duration": "24h"}', type: 'push', is_read: false, sent: false, sent_at: null, created_at: '2023-10-26T09:45:00Z' },
        { id: 5, user_id: 5, title: 'Emergency Test', body: 'This is a test of the emergency notification system', data: '{"type": "test", "system": "emergency"}', type: 'both', is_read: true, sent: true, sent_at: '2023-10-26T11:00:00Z', created_at: '2023-10-25T20:10:00Z' },
        { id: 6, user_id: 1, title: 'App Update', body: 'New features available in the latest update. Please update now.', data: '{"type": "app_update", "version": "2.1.0"}', type: 'in_app', is_read: false, sent: false, sent_at: null, created_at: '2023-10-26T14:20:00Z' }
    ]
};

// Hàm để lấy dữ liệu mẫu
window.getMockData = () => window.mockData;