window.onload = async function () 
{
    let canvas = document.getElementById("canvas");
    let scroll = document.getElementById("scroll")
    let context = canvas.getContext("2d");
    let quotes = document.querySelectorAll(".quote");
    // let quoteContainers = document.querySelectorAll(".quoteContainer");
    let search = document.querySelector("#search");

    context.font = window.getComputedStyle(scroll).fontSize.concat(" ").concat(window.getComputedStyle(scroll).fontFamily)

    let tildeWidth = context.measureText("_").width;
    let spaceWidth = context.measureText(" ").width;
    let pipeWidth = context.measureText("|").width;
    let pipeHeight = Number(window.getComputedStyle(scroll).fontSize.match(/\d+/)[0]) * 1.3;


    let widthMargin = 10; // this is measured in tilde-widths
    let heightMargin = 2; // this is measured in pipe-heights

    let shadowWizardMoneyGang = ["Consult the Enigmatic Knowledge...", "Peruse the Ancient Texts...", "Explore the Oracular Memorandums..."];
    search.textContent = shadowWizardMoneyGang[Math.floor(Math.random() * shadowWizardMoneyGang.length)];

    function quoteSize(quote) {
        quote.style.fontSize = String(Math.min(Math.max(3.5 - (0.06 * Math.abs(((quote.getBoundingClientRect().y / window.innerHeight * 100) - 47))), 2), 3)).concat("vmin");
    }
    quotes.forEach(quote => {
        quoteSize(quote);
    });

    async function doAnim() {
        for (i = 0; i <= Math.floor((window.innerHeight / pipeHeight)) - heightMargin - 5; i++) {
        let pageWidth = window.innerWidth;
        let top = "     ".concat("_".repeat(Math.floor(pageWidth / tildeWidth) - widthMargin - 7).concat("\n"));
        let section = "       |".concat(" ".repeat(Math.floor((pageWidth - (2 * pipeWidth + (widthMargin + 8) * tildeWidth)) / spaceWidth + 3))).concat("|").concat("\n");
        scrollText = top
            .concat("  /  \\".concat(" ".repeat(Math.floor((pageWidth - (-2 * spaceWidth + (1 + widthMargin + 8) * tildeWidth)) / spaceWidth))).concat("\\   \n"))
            .concat("|     |".concat(" ".repeat(Math.floor((pageWidth - (-1 * spaceWidth + (1 + widthMargin + 8) * tildeWidth)) / spaceWidth))).concat("|   \n"))
            .concat("  \\_ |".concat(" ".repeat(Math.floor((pageWidth - (-1 * spaceWidth + (1 + widthMargin + 8) * tildeWidth)) / spaceWidth))).concat("|   \n"))
            .concat(section.repeat(i))
            .concat("       |       ".concat("_".repeat(Math.floor((pageWidth - (2 * pipeWidth + (widthMargin + 8) * tildeWidth)) / tildeWidth - 1.5))).concat("|___").concat("\n"))
            .concat("       |     /".concat(" ".repeat(Math.floor((pageWidth - (2 * pipeWidth + (widthMargin + 8) * tildeWidth)) / spaceWidth - 1.5))).concat("     /").concat("\n"))
            .concat("        \\_/".concat("_".repeat(Math.floor((pageWidth - (2 * pipeWidth + (widthMargin + 8) * tildeWidth)) / tildeWidth + 1))).concat("_/").concat("\n"));
        scroll.textContent = scrollText;
        await new Promise(r => setTimeout(r, 15));
        }
    }
    await doAnim();

    function calcScroll() {
        let pageWidth = window.innerWidth;
        let pageHeight = window.innerHeight;
        let top = "     ".concat("_".repeat(Math.floor(pageWidth / tildeWidth) - widthMargin - 7).concat("\n"));
        let section = "       |".concat(" ".repeat(Math.floor((pageWidth - (2 * pipeWidth + (widthMargin + 8) * tildeWidth)) / spaceWidth + 3))).concat("|").concat("\n");
        scrollText = top
            .concat("  /  \\".concat(" ".repeat(Math.floor((pageWidth - (-2 * spaceWidth + (1 + widthMargin + 8) * tildeWidth)) / spaceWidth))).concat("\\   \n"))
            .concat("|     |".concat(" ".repeat(Math.floor((pageWidth - (-1 * spaceWidth + (1 + widthMargin + 8) * tildeWidth)) / spaceWidth))).concat("|   \n"))
            .concat("  \\_ |".concat(" ".repeat(Math.floor((pageWidth - (-1 * spaceWidth + (1 + widthMargin + 8) * tildeWidth)) / spaceWidth))).concat("|   \n"))
            .concat(section.repeat(Math.floor((pageHeight / pipeHeight)) - heightMargin - 5))
            .concat("       |       ".concat("_".repeat(Math.floor((pageWidth - (2 * pipeWidth + (widthMargin + 8) * tildeWidth)) / tildeWidth - 1.5))).concat("|___").concat("\n"))
            .concat("       |     /".concat(" ".repeat(Math.floor((pageWidth - (2 * pipeWidth + (widthMargin + 8) * tildeWidth)) / spaceWidth - 1.5))).concat("     /").concat("\n"))
            .concat("        \\_/".concat("_".repeat(Math.floor((pageWidth - (2 * pipeWidth + (widthMargin + 8) * tildeWidth)) / tildeWidth + 1))).concat("_/").concat("\n"));
        scroll.textContent = scrollText;
    }
    calcScroll();
    window.onresize = calcScroll;
    
    window.onscroll = function () {
        quotes.forEach(quote => {
            quoteSize(quote);
        });
    }

    // document.onscrollend = (event) => {
    //     let closest = quoteContainers[0];
    //     // quotes.forEach(quote => {
    //     //     if (Math.abs((quote.getBoundingClientRect().y + (quote.getBoundingClientRect().height / 2) - (window.innerHeight / 2))) < Math.abs((closest.getBoundingClientRect().y + (closest.getBoundingClientRect().height / 2) - (window.innerHeight / 2)))) {
    //     //         closest = quote;
    //     //     }
    //     // });

    //     quoteContainers.forEach(quote => {
    //         if (Math.abs((quote.getBoundingClientRect().top - (window.innerHeight / 2))) < Math.abs((closest.getBoundingClientRect().top - (window.innerHeight / 2)))) {
    //             closest = quote;
    //         }
    //     });
    //     // window.scrollTo({top: closest.getBoundingClientRect().y + (closest.getBoundingClientRect().height / 2) + document.body.scrollTop, behavior: "smooth"});
    //     console.log(closest.textContent)
    //     window.scrollTo({top: closest.getBoundingClientRect().y + document.body.scrollTop, behavior: "smooth"});
    //     // window.scrollBy({top: (window.innerHeight / 2) - closest.getBoundingClientRect().y + (closest.getBoundingClientRect().height / 2), behavior: "smooth"});
    //     // window.scrollBy({top: -((window.innerHeight / 2) - closest.getBoundingClientRect().y + (closest.getBoundingClientRect().height / 2)), behavior: "smooth"});
    // };
}