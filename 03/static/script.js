async function sendQuestion() {
    const questionInput = document.getElementById("question");
    const answerElement = document.getElementById("answer");

    if (!questionInput.value.trim()) {
        answerElement.innerText = "Vui lòng nhập câu hỏi!";
        return;
    }

    answerElement.innerText = "Đang xử lý...";

    try {
        const response = await fetch("/groq", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ prompt: questionInput.value.trim() })
        });

        if (!response.ok) {
            throw new Error(`Server trả về lỗi: ${response.status}`);
        }

        const data = await response.json();

        if (data.answer) {
            answerElement.innerText = data.answer;
        } else {
            answerElement.innerText = "Không nhận được phản hồi hợp lệ từ API.";
        }
    } catch (error) {
        answerElement.innerText = `Lỗi: ${error.message}`;
    }
}
