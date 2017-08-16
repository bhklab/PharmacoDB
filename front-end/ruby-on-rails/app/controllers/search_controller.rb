class SearchController < ApplicationController
    before_action :search_engine

    def index
    end

    private

    def search_engine
        @exit_code = 0

        if params[:q].present?
            terms = []
            input = params[:q]

            # if input == "\\markjoke"
            #     connection = open("http://deepmark.herokuapp.com")
            #     @text = JSON.parse(connection.read)["text"].split("\n").join("<br>")
            #     connection.close()
            #     @exit_code = 999
            #     return
            # end

            if ['cell', 'cells', 'cell line', 'cell lines'].include? input
                redirect_to "/cell_lines"
                return
            elsif ['tissue', 'tissues'].include? input
                redirect_to "/tissues"
                return
            elsif ['drug', 'drugs'].include? input
                redirect_to "/drugs"
                return
            elsif ['target', 'targets'].include? input
                redirect_to "/targets"
                return
            elsif ['dataset', 'datasets'].include? input
                redirect_to "/datasets"
                return
            elsif ['help'].include? input
                redirect_to "/docs"
            end

            type = ["cell", "tissue", "drug", "dataset", "source", "target"]
            paths = ["cell_lines", "tissues", "drugs", "datasets", "sources", "targets"]
            filter_type = [["SourceCellName", "SourceTissueName", "SourceDrugName"], ["Cell", "Tissue", "Drug", "Dataset", "Source", "Target"]]
            filter_name = ["cell_name", "tissue_name", "drug_name", "dataset_name", "source_name", "target_name"]
            filter_id = ["cell_id", "tissue_id", "drug_id", "dataset_id", "source_id", "target_id"]

            continue_loop = true
            filter_type.each do |f|
                unless continue_loop
                    break
                end
                max = f.count - 1
                for i in 0..max
                    object = eval(f[i]).where("#{filter_name[i]} = ?", input).first
                    if object.present?
                        continue_loop = false
                        terms << [object, i]
                        break
                    end
                end
            end

            if continue_loop
                input = params[:q].split(Rails.application.config.my_search_split)
                input.each do |p|
                    continue_loop = true
                    filter_type.each do |f|
                        unless continue_loop
                            break
                        end
                        max = f.count - 1
                        for i in 0..max
                            object = eval(f[i]).where("#{filter_name[i]} = ?", p).first
                            if object.present?
                                continue_loop = false
                                terms << [object, i]
                                break
                            end
                        end
                    end
                    if continue_loop
                        @error_message = ["Your search", "#{p}", "did not match any records in database."]
                        @error_datasets = Dataset.all.pluck(:dataset_name).join(", ")
                        @suggestions = ["Check for any spelling errors.", "Make sure to only search <a href=\"/cell_lines\">cell line</a>/<a href=\"/tissues\">tissue</a>/<a href=\"/drugs\">drug</a>/<a href=\"/datasets\">dataset</a> name(s).", "If searching multiple names, check the <a href=\"/docs\"> documentation </a> to see which combinations PharmacoDB handles."]
                        @exit_code = 1
                        return
                    end
                end
            end

            if terms.count == 1
                index = terms[0][1]
                object_id = terms[0][0].send(filter_id[index])

                redirect_to "/#{paths[index]}/#{object_id}"
                return
            end # /if terms.count == 1

            # XXX =================== update soon ==================================================

            ids_only = terms.map {|x| x[1]}

            if ids_only.uniq.count == 1
                index = terms[0][1]
                if index == 3
                    cell_line_ids = []
                    tissue_ids = []
                    drug_ids = []
                    round = 0
                    @names = ""
                    @search_quers = []

                    @cell_line_exps = []
                    @tissue_exps = []
                    @drug_exps = []
                    if type[index] == "dataset"
                        terms.each do |t|
                            # @exps.push(Experiment.where(dataset_id: t).includes(:cell).pluck(:cell_name).uniq)
                            @cell_line_exps.push(Experiment.where(dataset_id: t[0].send(filter_id[t[1]])).includes(:cell).pluck(:cell_name).uniq)
                            @tissue_exps.push(Experiment.where(dataset_id: t[0].send(filter_id[t[1]])).includes(:tissue).pluck(:tissue_name).uniq)
                            @drug_exps.push(Experiment.where(dataset_id: t[0].send(filter_id[t[1]])).includes(:drug).pluck(:drug_name).uniq)
                            @search_quers.push(t[0].send(filter_name[index]))
                            tempexp = Experiment.where(dataset_id: t[0].send(filter_id[t[1]]))
                            cid = tempexp.pluck(:cell_id)
                            tid = CellTissue.where(cell_id: cid).pluck(:tissue_id)
                            # tid = Tissue.joins(:cell_tissues).where(cell_tissues: {cell_id: cid}).pluck(:tissue_id).uniq
                            did = tempexp.pluck(:drug_id)
                            if round == 0
                                cell_line_ids = cid
                                tissue_ids = tid
                                drug_ids = did
                                @names = t[0].send(filter_name[index])
                                round += 1
                            else
                                cell_line_ids = cell_line_ids & cid
                                tissue_ids = tissue_ids & tid
                                drug_ids = drug_ids & did
                                @names = @names + ", " + t[0].send(filter_name[index])
                            end
                        end
                    end
                    @middle = []

                    cell_lines = Cell.where(cell_id: cell_line_ids).pluck(:cell_id, :cell_name)
                    venn_cell_lines = Cell.where(cell_id: cell_line_ids).pluck(:cell_name)
                    @cell_lines = cell_lines.paginate(:page => params[:page], :per_page => 10)
                    @cell_lines_count = cell_lines.count
                    tissues = Tissue.where(tissue_id: tissue_ids).pluck(:tissue_id, :tissue_name)
                    venn_tissues = Tissue.where(tissue_id: tissue_ids).pluck(:tissue_name)
                    @tissues = tissues.paginate(:page => params[:page], :per_page => 10)
                    @tissues_count = tissues.count
                    drugs = Drug.where(drug_id: drug_ids).pluck(:drug_id, :drug_name)
                    venn_drugs = Drug.where(drug_id: drug_ids).pluck(:drug_name)
                    @drugs = drugs.paginate(:page => params[:page], :per_page => 20)
                    @drugs_count = drugs.count
                    @middle.push(venn_cell_lines, venn_tissues, venn_drugs)

                    # require 'rinruby'

                    # R.image_path = Rails.root.join("app", "assets", "images", "venn.png").to_s
                    # R.lib_path = Rails.root.join("lib", "Rlibraries").to_s
                    # R.eval("library(rJava, lib.loc=lib_path)")
                    # R.eval("library(venneuler, lib.loc=lib_path)")
                    # R.eval("v <- venneuler(c(A=0.3, B=0.3, C=1.1, 'A&B'=0.1, 'A&C'=0.2, 'B&C'=0.1 ,'A&B&C'=0.1))")
                    # R.eval("png(filename=image_path, bg='transparent')")
                    # R.eval("plot(v)")
                    # R.eval("dev.off()")

                    @exit_code = 3
                    return
                end
            end

            if (terms.uniq.count == 2 || terms.uniq.count == 3) && (terms.count == 2 || terms.count == 3)
                if terms.count == 2 && (terms[0][1] == 0 && terms[1][1] == 2 || terms[1][1] == 0 && terms[0][1] == 2)
                    index1 = terms[0][1]
                    index2 = terms[1][1]


                    object1 = eval(filter_type[1][index1]).where("#{filter_id[index1]} = ?", terms[0][0].send(filter_id[index1])).first
                    object2 = eval(filter_type[1][index2]).where("#{filter_id[index2]} = ?", terms[1][0].send(filter_id[index2])).first
                    object1_name = object1.send(filter_name[index1])
                    object1_id = object1.send(filter_id[index1])
                    object2_name = object2.send(filter_name[index2])
                    object2_id = object2.send(filter_id[index2])


                    if index1 == 0
                        @cell_name = object1_name
                        @drug_name = object2_name
                    else
                        @drug_name = object1_name
                        @cell_name = object2_name
                    end

                    experiment = Experiment.where("#{filter_id[index1]} = ?", object1_id).where("#{filter_id[index2]} = ?", object2_id)

                    @highlighted_datasets = experiment.map{|x| x.dataset.dataset_name}.uniq

                elsif terms.count == 3 && (terms.pluck(1).to_set == [0,2,3].to_set)
                    # index1 = terms[0][1]
                    # index2 = terms[1][1]
                    # index3 = terms[2][1]

                    cell_term = terms[terms.pluck(1).each_index.detect{|x| terms.pluck(1)[x] == 0}]
                    drug_term = terms[terms.pluck(1).each_index.detect{|x| terms.pluck(1)[x] == 2}]
                    dataset_term = terms[terms.pluck(1).each_index.detect{|x| terms.pluck(1)[x] == 3}]

                    # object1 = eval(filter_type[1][index1]).where("#{filter_id[index1]} = ?", terms[0][0].send(filter_id[index1])).first
                    # object2 = eval(filter_type[1][index2]).where("#{filter_id[index2]} = ?", terms[1][0].send(filter_id[index2])).first
                    # object3 = eval(filter_type[1][index3]).where("#{filter_id[index3]} = ?", terms[2][0].send(filter_id[index3])).first
                    # object1_name = object1.send(filter_name[index1])
                    # object1_id = object1.send(filter_id[index1])
                    # object2_name = object2.send(filter_name[index2])
                    # object2_id = object2.send(filter_id[index2])
                    # object3_name = object3.send(filter_name[index3])
                    # object3_id = object3.send(filter_id[index3])


                    # if index1 == 0
                    #     @cell_name = object1_name
                    #     @drug_name = object2_name
                    # else
                    #     @drug_name = object1_name
                    #     @cell_name = object2_name
                    # end

                    experiment = Experiment.where("cell_id = ?", cell_term[0].cell_id).where("drug_id = ?", drug_term[0].drug_id)
                    @highlighted_datasets = [dataset_term[0].dataset_name]


                end
                @dataset_names = []
                @dataset_click_ids = []

                @profiles = []
                if experiment.present?
                    i = 0
                    @js_code = ""
                    @data = '{"data": ['
                    # dataset_name = ""
                    ii = 0
                    experiment.each do |e|
                        ii += 1
                        @dataset_names << e.dataset.dataset_name
                        @dataset_click_ids << e.dataset.dataset_name
                        @profiles << e.profile.attributes
                        dose = DoseResponse.where(experiment_id: e.experiment_id)
                        if dose.present?
                            # trace_name = "trace" + "#{i}"
                            if i == 0
                                @data = @data + "{"
                            else
                                @data = @data + ", " + "{"
                            end
                            i += 1
                            @data = @data + '"experiment_id": ' + e.id.to_s + ', "cell_line": { "id":' + e.cell.cell_id.to_s + ', "name": "' + e.cell.cell_name + '"}, ' + '"tissue": { "id":' + e.tissue_id.to_s + ', "name": "' + e.cell.tissues[0].tissue_name + '"}, ' + '"drug": { "id":' + e.drug.drug_id.to_s + ', "name": "' + e.drug.drug_name + '"}, ' + '"dataset": { "id":' + e.dataset.dataset_id.to_s + ', "name": "' + e.dataset.dataset_name + '"}, "params": {' + '"HS": ' + e.profile.HS.to_s +  ', "Einf": ' + e.profile.Einf.to_s + ', "EC50": ' + e.profile.EC50.to_s + '},' + ' "dose_responses":['

                            dose = dose.pluck(:dose)
                            # doses = "[" + dose.join(", ") + "]"
                            response = DoseResponse.where(experiment_id: e.experiment_id).pluck(:response)
                            # @range_max = response.max
                            # responses = "[" + response.join(", ") + "]"
                            # dataset_name = Dataset.where(dataset_id: e.dataset_id).first.dataset_name
                            for i in 0..dose.length-2
                                @data = @data + '{ "dose":' + dose[i].to_s + ', "response":' + response[i].to_s + '},'
                            end
                            @data = @data + '{ "dose":' + dose[dose.length-1].to_s + ', "response":' + response[dose.length-1].to_s + '} ]'


                            # @js_code = @js_code +
                            # "var #{trace_name} = {
                            # x: #{doses},
                            # y: #{responses},
                            # type: 'lines+markers',
                            # name: '#{dataset_name}'
                            # };\n"
                        end
                        if ii == experiment.length
                            @data = @data + '}'
                        else
                            @data = @data + '} '
                        end
                    end
                    @data = @data + ']}'
                    # @data = @data + "]"
                    # if i == 1
                    #     @title = "Dose/Response curve for #{object1_name} vs. #{object2_name} - #{dataset_name}"
                    # else
                    @title = "Dose/Response curve for #{object1_name} vs. #{object2_name}"
                    # end
                    @exit_code = 2

                    if !@dataset_names.empty?
                        dups = @dataset_names.group_by{|e| e}.map{|k, v| {k => v.length}}.find_all{|x| x.values.first>1}
                        unless dups.empty?
                            ddd = dups.map{|x| x.keys}.pluck(0)
                            for d in ddd
                                i = 1
                                for j in 0...@dataset_names.length
                                    if d.include? @dataset_names[j]
                                        @dataset_names[j] = @dataset_names[j] + " rep " + i.to_s
                                        @dataset_click_ids[j] = @dataset_click_ids[j] + "rep" + i.to_s

                                        i += 1
                                    end
                                end
                            end
                        end
                    end


                    return
                else
                    @error_message = ["This cell line to drug combination - ", "#{object1_name} vs. #{object2_name} ", "- has not been tested in our database."]
                    @error_datasets = Dataset.all.pluck(:dataset_name).join(", ")
                    @suggestions = ["Check for spelling errors.", "Try a different cell line to drug combination."]
                    @exit_code = 1
                    return
                end
            # end # /cell vs. drug, 1 to 1
                # end # /if terms.count == 2
            end # /if terms.uniq.count == 2

            @error_message = ["Your search - ", "#{params[:q]} ", "- did not yield any results."]
            @error_datasets = Dataset.all.pluck(:dataset_name).join(", ")
            @suggestions = ["Check for spelling errors.", "Try a different combination (see documentation for handled options)."]
            @exit_code = 1
            return

            # =================== end update ===================================================

        end #/ if params[:q].present?
        redirect_to '/'
    end #/ search_engine
end
