# TechmasterProject

Giới thiệu dự án: Đây là dự án từ công ty Techmaster dùng làm đề thi cho thực tập sinh golang,
dự án được chia ra làm 3 folder 01, 02, 03 tương ứng vs 3 bài

## Folder 01: Gọi vào Groq API
- Cấu trúc thư mục bao gồm:
* folder config chứa file config.go có chức năng LoadEnv dùng để tải các biến môi trường từ tệp .env vào ứng dụng
* folder handler chứa file handler.go có chức năng (xử lý yêu cầu) cho ứng dụng, được triển khai bằng Iris framework. 
Các handler này sẽ xử lý các yêu cầu HTTP từ người dùng và trả về phản hồi phù hợp. 
Trong file này, chúng ta có hai handler chính: ServeHomePage và HandleAsk.
* folder service chứa file groq.go có một hàm quan trọng có nhiệm vụ giao tiếp với Groq API để gửi yêu cầu và nhận câu trả lời. 
Hàm này được gọi trong các phần khác của ứng dụng để lấy câu trả lời từ dịch vụ AI của Groq.
* folder template chứa giao diện có đuôi .html
* file main
- Cách chạy:
+ Tạo file .env bằng cách vào đường link "https://console.groq.com/keys"
+ Login tài khoản, ấn Create API Key và sao chép API Key dán vào file .env có dạng như sau:
  GROQ_API_KEY=api_key_vua_sao_chep
+ Mở rồi chạy file index.html trong thư mục template

## Folder 02: Sinh file SSML từ hội thoại
- Bao gồm 2 file HTML và Javascript như đề bài yêu cầu
- Cách chạy mở rồi chạy file index.html
- Vì gấp gáp nên em chưa kịp CSS cho giao diện mong thầy thông cảm, có thời gian em sẽ quay lại làm

## Folder 03: Tạo hội thoại từ prompt và trích xuất từ mới trong hội thoại
- Cấu trúc thư mục bao gồm:
+ folder db chứa file database.go có hàm ConnectToDB làm nhiệm vụ kết nối tới postgresql
+ folder processing chứa file text_processing.go có chức năng xử lý chuỗi(lọc chuỗi, chuyển đổi chuỗi sang json, tạo csdl, lưu vào csdl, ...)
+ file conversation.txt chứa đoạn hội thoại được sinh sau khi chạy file main.go
+ file docker-compose.yaml để chạy container kết nối đên csdl postgresql
+ file main
- Cách chạy: 
+ tải và bật docker lên sau đó chạy câu lệnh "docker-compose up -d"
+ cd (điều hướng) đến thư mục 03, trong terminal gõ "go run main.go"
+ tại terminal chạy lệnh "docker exec -it my_postgres psql -U myuser -d mydatabase"
    sau đó gõ "\dt" để thấy được những bảng đã tạo
            gõ "\dt dialog" để xem được chi tiết bảng
    Thêm bảng word_dialog bằng cách copy
  "CREATE TABLE word_dialog (
  dialog_id BIGINT REFERENCES dialog(id) ON DELETE CASCADE,
  word_id BIGINT REFERENCES word(id) ON DELETE CASCADE,
  PRIMARY KEY (dialog_id, word_id)
  );" và dán vào terminal
    sau đó gõ \q để thoát ra
+ Chạy lại "go run main.go"
- Cũng vì gấp gáp nên em chưa làm phần tự động hóa bằng giao diện sử dụng framework Iris,
có thời gian em sẽ quay lại tìm hiểu và code tiếp