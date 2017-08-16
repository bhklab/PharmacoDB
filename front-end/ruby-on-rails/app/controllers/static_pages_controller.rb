class StaticPagesController < ApplicationController

	def home
		if not cookies.key?(:visited)
			cookies[:visited] = {
				:value => 0,
				:expires => 1.year.from_now
			}
		else
			# cookies[:visited] = {
			# 	:value => true,
			# 	:expires => 1.year.from_now
			# }
		end

		@cell_lines_count = Cell.count
		@tissues_count = Tissue.count
		@drugs_count = Drug.count
		@datasets_count = Dataset.count
		@experiments_count = Experiment.count
		@targets_count = Target.count

		if params[:title].present? && params[:message].present?
			@title = params[:title]
			@message = params[:message]

			@name = params[:name].present? ? params[:name] : "n/a"
			@email = params[:email].present? ? params[:email] : "n/a"
			@subject = params[:subject].present? ? params[:subject] : "n/a"

			Github.configure do |c|
				 c.basic_auth = 'user:pw'
				 c.user       = 'bhklab'
				 c.repo       = 'pharmacoDB'
			end
			github = Github.new
			github.issues.create title: @title, body: @message, labels: [@subject]
		end
	end

	def cite_us

	end

	def pharmacogx
		@dataset_id = []
		@dataset_name = []
		if params[:pgx].present?
			dataset = Dataset.find(params[:pgx].to_i)
			@dataset_id << dataset.dataset_id
			@dataset_name << dataset.dataset_name
		end
	end

end
