function generateSSML() {
    let conversation = document.getElementById("conversation").value;
    let voiceA = document.getElementById("voiceA").value;
    let voiceB = document.getElementById("voiceB").value;
    let lines = conversation.split("\n");
    let ssml = "<speak xml:lang='vi-VN'>\n";

    lines.forEach((line, index) => {
        let voice = (index % 2 === 0) ? voiceA : voiceB;
        ssml += `    <voice name='${voice}'>${line}</voice>\n`;
    });

    ssml += "</speak>";
    document.getElementById("ssmlOutput").textContent = ssml;
}
