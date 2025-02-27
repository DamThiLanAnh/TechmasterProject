async function sendQuestion() {
    let question = document.getElementById("question").value;
    if (!question.trim()) return alert("Vui lòng nhập câu hỏi!");

    let response = await fetch("/ask", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ question }),
    });

    let data = await response.json();
    document.getElementById("answer").innerText = data.answer || "Lỗi khi nhận phản hồi!";
}
