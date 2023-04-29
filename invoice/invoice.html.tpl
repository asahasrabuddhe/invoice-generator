<html lang="en">
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
    <title>Invoice - {{ .InvoiceNumber }} | {{ .InvoiceDate }}</title>
</head>
<body>
<section class="white">
    <div class="w-a4 h-a4 mx-auto bg-white">
        <article class="overflow-hidden">
            <div class="bg-[white] rounded-b-md">
                <div class="p-8">
                    <div class="grid grid-flow-col">
                        <!-- name and email -->
                        <div class="space-y-2 text-slate-700">
                            <p class="font-body font-light">{{ .FromEmail }}</p>
                            <p class="text-xl font-extrabold tracking-tight uppercase font-body">
                                {{ .FromName }}
                            </p>
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
                            <p class="font-body">{{ .FromAddress.Line1 }}</p>
                            <p class="font-body">{{ .FromAddress.Line2 }}</p>
                            <p class="font-body">{{ .FromAddress.Line3 }}</p>
                        </div>
                        <!-- invoice number and date -->
                        <div class="text-slate-700 font-light">
                            <p class="font-body text-right"><span class="font-bold">INVOICE #</span> {{ .InvoiceNumber
                                }}</p>
                            <p class="font-body text-right"><span class="font-bold">{{ .InvoiceDate }}</span></p>
                        </div>
                    </div>
                </div>
                <div class="px-8 py-2">
                    <!-- contact info: phone and email -->
                    <div class="text-slate-700 font-light">
                        <p class="font-body">{{ .FromPhone1 }}</p>
                        <p class="font-body">{{ .FromPhone2 }}</p>
                        <p class="font-body py-2">{{ .FromEmail }}</p>
                    </div>
                </div>
                <div class="px-8 py-2">
                    <p class="font-body font-bold text-slate-700">To</p>
                </div>
                <div class="px-8 py-2">
                    <!-- contact info: phone and email -->
                    <div class="text-slate-700 font-light">
                        <p class="font-body py-2">{{ .ToName }}</p>
                        <p class="font-body">{{ .ToAddress.Line1 }}</p>
                        <p class="font-body">{{ .ToAddress.Line2 }}</p>
                        <p class="font-body">{{ .ToAddress.Line3 }}</p>
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
                            <tr class="border-b border-slate-200">
                                <td class="py-4 pl-4 pr-3 pl-6 w-3/4">
                                    <p class="text-sm text-left font-medium text-slate-700">{{ formatDescription .Description}}</p>
                                </td>
                                <td class="py-4 pl-3 pr-4 pr-6 w-1/4">
                                    <p class="text-sm text-right text-slate-500">{{ formatAmount .Amount }}</p>
                                </td>
                            </tr>
                            {{ end }}
                            </tbody>
                            <tfoot>
                            <!--                            <tr>-->
                            <!--                                <th scope="row" colspan="3"-->
                            <!--                                    class="hidden pt-6 pl-6 pr-3 text-sm font-light text-right text-slate-500 sm:table-cell md:pl-0">-->
                            <!--                                    Subtotal-->
                            <!--                                </th>-->
                            <!--                                <th scope="row"-->
                            <!--                                    class="pt-6 pl-4 pr-3 text-sm font-light text-left text-slate-500 sm:hidden">-->
                            <!--                                    Subtotal-->
                            <!--                                </th>-->
                            <!--                                <td class="pt-6 pl-3 pr-4 text-sm text-right text-slate-500 sm:pr-6 md:pr-0">-->
                            <!--                                    $0.00-->
                            <!--                                </td>-->
                            <!--                            </tr>-->
                            <!--                            <tr>-->
                            <!--                                <th scope="row" colspan="3"-->
                            <!--                                    class="hidden pt-4 pl-6 pr-3 text-sm font-light text-right text-slate-500 sm:table-cell md:pl-0">-->
                            <!--                                    Tax-->
                            <!--                                </th>-->
                            <!--                                <th scope="row"-->
                            <!--                                    class="pt-4 pl-4 pr-3 text-sm font-light text-left text-slate-500 sm:hidden">-->
                            <!--                                    Tax-->
                            <!--                                </th>-->
                            <!--                                <td class="pt-4 pl-3 pr-4 text-sm text-right text-slate-500 sm:pr-6 md:pr-0">-->
                            <!--                                    $0.00-->
                            <!--                                </td>-->
                            <!--                            </tr>-->
                            <tr>
                                <th scope="row"
                                    class="pt-4 pl-6 pr-3 text-sm font-normal text-right text-slate-700 table-cell">
                                    Total
                                </th>
                                <td class="pt-4 pl-3 pr-4 text-sm font-normal text-right text-slate-700 pr-6">
                                    {{ formatAmount .Total }}
                                </td>
                            </tr>
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