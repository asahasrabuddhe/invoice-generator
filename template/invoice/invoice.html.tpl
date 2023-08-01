<html lang="en" xmlns="http://www.w3.org/1999/html">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width,initial-scale=0.8">
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            theme: {
                fontSize: {
                    sm: '0.8rem',
                    base: '1rem',
                    xl: '1.25rem',
                    '2xl': '1.563rem',
                    '3xl': '1.953rem',
                    '4xl': '2.441rem',
                    '5xl': '3.052rem',
                },
                extend: {
                    width: {
                        'a4': '210mm',
                    },
                    height: {
                        'a4': '297mm',
                    },
                }
            }
        }
    </script>
    <style>
        html, body {
            width: fit-content !important;
            height: fit-content !important;
        }
    </style>

    <style id=page_style>
        @page {
            size: 210mm 298mm !important;
        }
    </style>
    <title>Invoice - {{ .Number }} | {{ .Date }}</title>
</head>
<body>
<section class="white">
    <div class="w-a4 h-a4 mx-auto bg-white">
        <article class="overflow-hidden">
            <div class="bg-[white] rounded-b-md">
                <div class="px-8 py-4">
                    <div class="grid grid-flow-col">
                        <!-- name and email -->
                        <div class="space-y-2 text-slate-700">
                            <p class="font-body font-light">{{ .From.Email }}</p>
                            <p class="text-xl font-extrabold tracking-tight uppercase font-body">
                                {{ .From.Name }}
                            </p>
                            {{ if ne .From.AccountNumber "" }}
                            <p class="text-xl font-extrabold tracking-tight uppercase font-body">
                                Account Number - {{ .From.AccountNumber }}
                            </p>
                            {{ end }}
                            {{ if ne .From.GSTIN "" }}
                            <p class="text-xl font-extrabold tracking-tight uppercase font-body">
                                GSTIN - {{ .From.GSTIN }}
                            </p>
                            {{ end }}
                        </div>
                        <!-- invoice -->
                        <div class="space-y-2 text-slate-900">
                            <p class="text-2xl text-right font-extrabold tracking-tight uppercase font-body">
                                INVOICE
                            </p>
                        </div>
                    </div>
                </div>
                <div class="px-8 py-2">
                    <div class="grid grid-flow-col">
                        <!-- address -->
                        <div class="text-slate-700 font-light">
                            {{ range .From.AddressLines }}
                            <p class="font-body">{{ . }}</p>
                            {{ end }}
                        </div>
                        <!-- invoice number and date -->
                        <div class="text-slate-700 font-light">
                            <p class="font-body text-right"><span
                                    class="font-bold">INVOICE #</span> {{ .Number }}</p>
                            <p class="font-body text-right"><span
                                    class="font-bold">{{ .Date }}</span></p>
                        </div>
                    </div>
                </div>
                <div class="px-8 py-2">
                    <!-- contact info: phone and email -->
                    <div class="text-slate-700 font-light">
                        {{ range .From.Phone }}
                        <p class="font-body">{{ . }}</p>
                        {{ end }}
                        <p class="font-body py-2">{{ .From.Email }}</p>
                    </div>
                </div>
                <div class="px-8 py-2">
                    <p class="font-body font-bold text-slate-700">To</p>
                </div>
                <div class="px-8 py-2">
                    <!-- contact info: phone and email -->
                    <div class="text-slate-700 font-light">
                        <p class="font-body py-2">{{ .To.Name }}</p>
                        {{ if ne .To.GSTIN "" }}
                        <p class="font-body py-2">
                            GSTIN - {{ .To.GSTIN }}
                        </p>
                        {{ end }}
                        {{ range .To.AddressLines }}
                        <p class="font-body">{{ . }}</p>
                        {{ end }}
                    </div>
                </div>
                <div class="px-8 py-2">
                    <div class="flex flex-col mx-0 mt-8">
                        <table class="min-w-full divide-y divide-slate-500 table-fixed">
                            <thead>
                            <tr>
                                <th scope="col"
                                    class="py-3.5 pl-4 pr-3 text-left text-sm font-normal text-slate-700 pl-6 w-3/4"
                                >
                                    Description
                                </th>
                                <th scope="col"
                                    class="py-3.5 pl-3 pr-4 text-right text-sm font-normal text-slate-700 pr-6 w-1/4">
                                    Amount
                                </th>
                            </tr>
                            </thead>
                            <tbody>
                            {{ range .Lines }}
                                {{ if ne .Amount 0.0 }}
                                <tr class="border-b border-slate-200">
                                    <td class="py-4 pl-4 pr-3 pl-6 w-3/4">
                                        <p class="text-sm text-left font-medium text-slate-700">{{
                                            formatDescription
                                            .Description}}</p>
                                    </td>
                                    <td class="py-4 pl-3 pr-4 pr-6 w-1/4">
                                        <p class="text-sm text-right text-slate-700 font-bold">{{
                                            formatAmount $.Currency .Amount }}</p>
                                    </td>
                                </tr>
                                {{ end }}
                            {{ end }}
                            </tbody>
                            <tfoot>
                            {{ range $i, $t := .Tax }}
                                {{ if eq $i 0 }}
                                <tr>
                                    <th scope="row"
                                        class="pt-4 pl-6 pr-3 text-sm font-bold text-right text-slate-700 table-cell">
                                        Subtotal
                                    </th>
                                    <td class="pt-4 pl-3 pr-4 text-sm font-bold text-right text-slate-700 pr-6">
                                        {{ formatAmount $.Currency $.Total }}
                                    </td>
                                </tr>
                                {{ end }}
                                <tr>
                                    {{ if ne $t.AccountNumber "" }}
                                    <th scope="row"
                                        class="pt-4 pl-6 pr-3 text-sm font-bold text-right text-slate-700 table-cell">
                                        {{ $t.Type }} {{ $t.Rate }}% ({{ $t.AccountNumber }})
                                    </th>
                                    {{ else }}
                                    <th scope="row"
                                        class="pt-4 pl-6 pr-3 text-sm font-bold text-right text-slate-700 table-cell">
                                        {{ $t.Type }} {{ $t.Rate }}%
                                    </th>
                                    {{ end }}
                                    <td class="pt-4 pl-3 pr-4 text-sm font-bold text-right text-slate-700 pr-6">
                                        {{ formatAmount $.Currency ($t.Total $.Total) }}
                                    </td>
                                </tr>
                            {{ end }}
                            <tr>
                                <th scope="row"
                                    class="pt-4 pl-6 pr-3 text-sm font-bold text-right text-slate-700 table-cell">
                                    Total
                                </th>
                                <td class="pt-4 pl-3 pr-4 text-sm font-bold text-right text-slate-700 pr-6">
                                    {{ formatAmount $.Currency (add $.Total ($.Tax.Total $.Total)) }}
                                </td>
                            </tr>
                            {{ if ne .CurrencyRate 0.0 }}
                            <tr>
                                <th scope="row"
                                    class="pt-4 pl-6 pr-3 text-sm font-bold text-right text-slate-700 table-cell">
                                    1 {{ .Currency }} = {{ .CurrencyRate }} INR
                                </th>
                                <td class="pt-4 pl-3 pr-4 text-sm font-bold text-right text-slate-700 pr-6">
                                    {{ formatAmount .Currency (mul (add .Total (.Tax.Total .Total)) .CurrencyRate) }}
                                </td>
                            </tr>
                            {{ end }}
                            </tfoot>
                        </table>
                    </div>
                </div>
            </div>
        </article>
    </div>
</section>
</body>
</html>