<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html lang="en">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <title>Invoice - {{ .InvoiceNumber }} | {{ .InvoiceDate }}</title>
  <meta name="author" content="Ajitem Sahasrabuddhe">
  <style type="text/css">
      body {
          margin-top: 0;
          margin-left: 0;
      }

      #page_1 {
          position: relative;
          overflow: hidden;
          margin: 36px 0 194px 47px;
          padding: 0;
          border: none;
          width: 819px;
      }

      #page_1 #id1_1 {
          border: none;
          margin: 0 0 0 0;
          padding: 0;
          width: 781px;
          overflow: hidden;
      }

      #page_1 #id1_2 {
          border: none;
          margin: 54px 0 0 246px;
          padding: 0;
          width: 523px;
          overflow: hidden;
      }


      .ft0 {
          font: 15px 'Century Gothic';
          color: #404040;
          line-height: 20px;
      }

      .ft1 {
          font: 1px 'Calibri';
          line-height: 1px;
      }

      .ft2 {
          font: bold 35px 'Century Gothic';
          color: #2e74b5;
          line-height: 24px;
      }

      .ft3 {
          font: 1px 'Calibri';
          line-height: 3px;
      }

      .ft4 {
          font: bold 19px 'Century Gothic';
          color: #7f7f7f;
          line-height: 23px;
      }

      .ft5 {
          font: bold 15px 'Century Gothic';
          color: #2e74b5;
          line-height: 18px;
      }

      .ft6 {
          font: 15px 'Century Gothic';
          color: #404040;
          line-height: 19px;
      }

      .ft7 {
          font: 15px 'Century Gothic';
          color: #404040;
          line-height: 18px;
      }

      .ft8 {
          font: 15px 'Century Gothic';
          color: #404040;
          line-height: 17px;
      }

      .ft9 {
          font: 1px 'Calibri';
          line-height: 17px;
      }

      .ft10 {
          font: 15px 'Century Gothic';
          text-decoration: underline;
          color: #0563c1;
          line-height: 20px;
      }

      .ft11 {
          font: 1px 'Calibri';
          line-height: 4px;
      }

      .ft12 {
          font: 1px 'Calibri';
          line-height: 5px;
      }

      .ft14 {
          font: bold 15px 'Calibri';
          line-height: 18px;
      }

      .ft16 {
          font: 1px 'Calibri';
          line-height: 6px;
      }

      .ft17 {
          font: 1px 'Calibri';
          line-height: 7px;
      }

      .ft18 {
          font: bold 15px 'Century Gothic';
          color: #404040;
          line-height: 18px;
      }

      .ft19 {
          font: 15px 'Century Gothic';
          color: #2e74b5;
          line-height: 20px;
      }

      .p0 {
          text-align: left;
          margin-top: 0;
          margin-bottom: 0;
          /*white-space: nowrap;*/
      }

      .p1 {
          text-align: right;
          padding-right: 1px;
          margin-top: 0;
          margin-bottom: 0;
          white-space: nowrap;
      }

      .p2 {
          text-align: right;
          margin-top: 0;
          margin-bottom: 0;
          white-space: nowrap;
      }

      .p3 {
          text-align: right;
          padding-right: 167px;
          margin-top: 0;
          margin-bottom: 0;
          white-space: nowrap;
      }

      .p4 {
          text-align: right;
          padding-right: 166px;
          margin-top: 0;
          margin-bottom: 0;
          white-space: nowrap;
      }

      .p5 {
          text-align: left;
          padding-left: 1px;
          margin-top: 0;
          margin-bottom: 0;
          white-space: nowrap;
      }

      .p6 {
          text-align: left;
          margin-top: 0;
          margin-bottom: 0;
      }

      .td0 {
          padding: 0;
          margin: 0;
          width: 445px;
          vertical-align: bottom;
      }

      .td1 {
          padding: 0;
          margin: 0;
          width: 59px;
          vertical-align: bottom;
      }

      .td2 {
          padding: 0;
          margin: 0;
          width: 116px;
          vertical-align: bottom;
      }

      .td3 {
          padding: 0;
          margin: 0;
          width: 225px;
          vertical-align: top;
      }

      .td4 {
          border-bottom: #bdd6ee 1px solid;
          padding: 0;
          margin: 0;
          width: 495px;
          vertical-align: bottom;
      }

      .td5 {
          border-bottom: #bdd6ee 1px solid;
          padding: 0;
          margin: 0;
          width: 59px;
          vertical-align: bottom;
      }

      .td6 {
          border-bottom: #9cc2e5 1px solid;
          padding: 0;
          margin: 0;
          width: 495px;
          vertical-align: bottom;
      }

      .td7 {
          border-bottom: #9cc2e5 1px solid;
          padding: 0;
          margin: 0;
          width: 59px;
          vertical-align: bottom;
      }

      .td10 {
          padding: 0;
          margin: 0;
          width: 490px;
          vertical-align: bottom;
      }

      .td11 {
          padding: 0;
          margin: 0;
          width: 137px;
          vertical-align: bottom;
      }

      .td12 {
          border-bottom: #bdd6ee 1px solid;
          padding: 0;
          margin: 0;
          width: 490px;
          vertical-align: bottom;
      }

      .td13 {
          border-bottom: #bdd6ee 1px solid;
          padding: 0;
          margin: 0;
          width: 65px;
          vertical-align: bottom;
      }

      .tr0 {
          height: 32px;
      }

      .tr1 {
          height: 29px;
      }

      .tr2 {
          height: 3px;
      }

      .tr3 {
          height: 28px;
      }

      .tr4 {
          height: 46px;
      }

      .tr5 {
          height: 19px;
      }

      .tr6 {
          height: 18px;
      }

      .tr7 {
          height: 33px;
      }

      .tr8 {
          height: 17px;
      }

      .tr9 {
          height: 35px;
      }

      .tr10 {
          height: 36px;
      }

      .tr11 {
          height: 43px;
      }

      .tr12 {
          height: 34px;
      }

      .tr13 {
          height: 20px;
      }

      .tr14 {
          height: 4px;
      }

      .tr15 {
          height: 5px;
      }

      .tr16 {
          height: 21px;
      }

      .tr18 {
          height: 6px;
      }

      .tr19 {
          height: 7px;
      }

      .t0 {
          width: 770px;
          margin-left: 1px;
          font: bold 15px 'Calibri';
      }

      .t1 {
          width: 555px;
          font: bold 15px 'Century Gothic';
          color: #2e74b5;
      }

      table {
          table-layout: fixed;
      }

      table td {
          word-wrap: break-word;
      }

  </style>
  <link rel="stylesheet" href="font.css"/>
</head>

<body>
<div id="page_1">
  <div id="id1_1">
    <table cellpadding="0" cellspacing="0" class="t0">
      <tbody>
      <tr>
        <td rowspan="2" class="tr0 td0"><p class="p0 ft0">ajitem.s@outlook.com</p></td>
        <td class="tr1 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr1 td2"><p class="p1 ft2">INVOICE</p></td>
      </tr>
      <tr>
        <td class="tr2 td1"><p class="p0 ft3">&nbsp;</p></td>
        <td class="tr2 td2"><p class="p0 ft3">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr3 td0"><p class="p0 ft4">Ajitem Sahasrabuddhe</p></td>
        <td class="tr3 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr3 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr4 td0"><p class="p0 ft0">1701, Trendy Towers 30,</p></td>
        <td class="tr4 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr4 td2"><p class="p2 ft5">INVOICE # <span class="ft0">{{ .InvoiceNumber }}</span></p></td>
      </tr>
      <tr>
        <td class="tr5 td0"><p class="p0 ft6">Amanora Park Town,</p></td>
        <td colspan="2" class="tr5 td3">
          <p class="p2 ft5">
            <nobr>{{ .InvoiceDate }}</nobr>
          </p>
        </td>
      </tr>
      <tr>
        <td class="tr6 td0"><p class="p0 ft7">Hadapsar, Pune – 411028</p></td>
        <td class="tr6 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr6 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr7 td0"><p class="p0 ft0">+91 888 832 4979</p></td>
        <td class="tr7 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr7 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr8 td0"><p class="p0 ft8">+91 928 445 1979</p></td>
        <td class="tr8 td1"><p class="p0 ft9">&nbsp;</p></td>
        <td class="tr8 td2"><p class="p0 ft9">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr9 td0"><p class="p0 ft10"><a href="mailto:ajitem.s@outlook.com">ajitem.s@outlook.com</a></p></td>
        <td class="tr9 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr9 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr10 td0"><p class="p0 ft5">TO</p></td>
        <td class="tr10 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr10 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr11 td0"><p class="p0 ft0">engineering.com</p></td>
        <td class="tr11 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr11 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr12 td0"><p class="p0 ft0">5285 Solar Drive, Suite 101,</p></td>
        <td class="tr12 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr12 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr6 td0"><p class="p0 ft7">Mississauga, Ontario,</p></td>
        <td class="tr6 td1"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr6 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr8 td0"><p class="p0 ft8">L4W 5B8 CANADA</p></td>
        <td class="tr8 td1"><p class="p0 ft9">&nbsp;</p></td>
        <td class="tr8 td2"><p class="p0 ft9">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr6 td4"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr6 td5"><p class="p0 ft1">&nbsp;</p></td>
        <td class="tr5 td2"><p class="p0 ft1">&nbsp;</p></td>
      </tr>
      <tr>
        <td class="tr13 td0"><p class="p0 ft5">Description</p></td>
        <td colspan="2" class="tr13 td3"><p class="p3 ft5">Amount</p></td>
      </tr>
      <tr>
        <td class="tr14 td6"><p class="p0 ft11">&nbsp;</p></td>
        <td class="tr14 td7"><p class="p0 ft11">&nbsp;</p></td>
        <td class="tr15 td2"><p class="p0 ft12">&nbsp;</p></td>
      </tr>
      {{ range .Lines }}
      <tr>
        <td class="tr16 td0"><p class="p0 ft14">{{ formatDescription .Description }}</p></td>
        <td colspan="2" class="tr16 td3"><p class="p4 ft14" style="text-align: right">{{ .Amount }}</p></td>
      </tr>
      <tr>
        <td class="tr18 td4"><p class="p0 ft16">&nbsp;</p></td>
        <td class="tr18 td5"><p class="p0 ft16">&nbsp;</p></td>
        <td class="tr19 td2"><p class="p0 ft17">&nbsp;</p></td>
      </tr>
      {{ end }}
      </tbody>
    </table>
    <table cellpadding="0" cellspacing="0" class="t1">
      <tbody>
      <tr>
        <td class="tr16 td10"><p class="p5 ft5">Total</p></td>
        <td class="tr16 td11"><p class="p0 ft18" style="white-space: nowrap; text-align: right">USD {{ .Total }}</p></td>
      </tr>
      <tr>
        <td class="tr2 td12"><p class="p0 ft3">&nbsp;</p></td>
        <td class="tr2 td13"><p class="p0 ft3">&nbsp;</p></td>
      </tr>
      </tbody>
    </table>
  </div>
  <div id="id1_2">
    <p class="p6 ft19">THANK YOU FOR YOUR BUSINESS!</p>
  </div>
</div>
</body>
</html>
