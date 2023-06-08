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
    <title>Summary Sheet {{ .Month }}</title>
</head>
<body>
<section class="white">
    <div class="w-a4 h-a4 mx-auto bg-white">
        <article class="overflow-hidden">
            <div class="bg-[white] rounded-b-md">
                <div class="px-8 py-4">
                    <div class="grid grid-flow-col">
                        <!-- invoice -->
                        <div class="space-y-2 text-slate-900">
                            <p class="text-2xl text-center font-extrabold tracking-tight uppercase font-body">
                                Summary Statement
                            </p>
                            <p class="font-light text-center">{{ .Month }}</p>
                        </div>
                    </div>
                </div>
                <div class="px-8 py-2">
                    <div class="flex flex-col mx-0 mt-8">
                        <table class="min-w-full divide-y divide-slate-500 table-fixed">
                            <thead>
                            <tr>
                                <th scope="col"
                                    class="py-3.5 pr-3 text-left text-sm font-normal text-slate-700 pl-6 w-1/3">
                                    Resource
                                </th>
                                <th scope="col"
                                    class="py-3.5 pr-3 text-left text-sm font-normal text-slate-700 pl-6 w-1/3">
                                    Invoice #
                                </th>
                                <th scope="col"
                                    class="py-3.5 pr-4 text-left text-sm font-normal text-slate-700 pl-6 w-1/3">
                                    Amount (US$)
                                </th>
                            </tr>
                            </thead>
                            <tbody>
                            {{ range .Lines }}
                            <tr class="border-b border-slate-200">
                                <td class="py-4 pl-4 pr-3 pl-6 w-1/3">
                                    <p class="text-sm text-left font-medium text-slate-700">{{ .Resource }}</p>
                                </td>
                                <td class="py-4 pl-6 pr-4 pr-6 w-1/3">
                                    <p class="text-sm text-left text-slate-700 font-bold">{{ .InvoiceNumber }}</p>
                                </td>
                                <td class="py-4 pl-6 pr-4 pr-6 w-1/3">
                                    <p class="text-sm text-left text-slate-700 font-bold">{{ formatAmount .Amount }}</p>
                                </td>
                            </tr>
                            {{ end }}
                            </tbody>
                            <tfoot>
                            <tr>
                                <th scope="row"
                                    colspan="2"
                                    class="pt-4 pl-6 pr-3 text-sm font-bold text-right text-slate-700 table-cell">
                                    Total
                                </th>
                                <td class="pt-4 pl-6 pr-4 text-sm font-bold text-left text-slate-700 pr-6">
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